package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"portal/internal/auth"
	"portal/internal/config"
	"portal/internal/handler"
	"portal/internal/kcadmin"
	"portal/internal/middleware"
	"portal/internal/permission"
	"portal/internal/repository"
	"portal/internal/service"
	sessionpkg "portal/internal/session"
	syncsvc "portal/internal/sync"
)

func main() {
	cfg := config.MustLoad()
	logger := config.NewLogger(cfg.Log.Level)
	ctx := context.Background()

	db, err := repository.NewMongo(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to connect to mongo", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if closeErr := db.Close(context.Background()); closeErr != nil {
			logger.Error("failed to close mongo", slog.Any("error", closeErr))
		}
	}()

	repos := repository.NewRepositories(db.Database, logger)
	oidcClient, err := auth.NewOIDCClient(ctx, cfg)
	if err != nil {
		logger.Error("failed to initialize oidc", slog.Any("error", err))
		os.Exit(1)
	}

	kcClient := kcadmin.NewClient(cfg, logger)
	syncService := syncsvc.NewService(kcClient, repos, logger)
	permissionService := permission.NewService(repos, cfg)
	sessionManager := sessionpkg.NewManager(repos.Sessions, logger, cfg)
	appService := service.NewAppService(permissionService, repos, cfg)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS(cfg))

	healthHandler := handler.NewHealthHandler(db, oidcClient)
	authHandler := handler.NewAuthHandler(cfg, oidcClient, syncService, sessionManager, repos, permissionService, logger)
	appHandler := handler.NewAppHandler(appService)
	adminHandler := handler.NewAdminHandler(repos, cfg)

	router.GET("/healthz", healthHandler.Healthz)
	router.GET("/readyz", healthHandler.Readyz)
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.File(cfg.Server.OpenAPIFilePath)
	})

	api := router.Group("/api/v1")
	{
		api.GET("/auth/login", authHandler.Login)
		api.GET("/auth/callback", authHandler.Callback)
		api.GET("/auth/logout", authHandler.Logout)

		protected := api.Group("/")
		protected.Use(middleware.Session(sessionManager))
		protected.Use(middleware.IdleTimeout(sessionManager))
		{
			protected.GET("/me", appHandler.Me)
			protected.GET("/apps", appHandler.Apps)

			admin := protected.Group("/admin")
			admin.Use(middleware.RequireAdmin(cfg))
			{
				admin.GET("/client-metas", adminHandler.ListClientMetas)
				admin.POST("/client-metas", adminHandler.UpsertClientMeta)
				admin.PUT("/client-metas/:clientId", adminHandler.UpsertClientMeta)
				admin.DELETE("/client-metas/:clientId", adminHandler.DeleteClientMeta)
				admin.GET("/settings", adminHandler.GetSettings)
				admin.PUT("/settings", adminHandler.UpdateSettings)
			}
		}
	}

	server := &http.Server{
		Addr:              cfg.Server.ListenAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("portal-api listening", slog.String("addr", cfg.Server.ListenAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server stopped unexpectedly", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("portal-api shutdown complete")
}
