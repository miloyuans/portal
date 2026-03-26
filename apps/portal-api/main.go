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

	adminClient := kcadmin.NewClient(cfg)
	syncService := syncsvc.NewService(adminClient, repos, cfg, logger)
	permissionService := permission.NewService(repos)
	sessionManager := sessionpkg.NewManager(repos.Sessions, cfg)
	appService := service.NewAppService(permissionService, repos)

	authHandler := handler.NewAuthHandler(cfg, oidcClient, syncService, sessionManager, repos, logger)
	portalHandler := handler.NewPortalHandler(appService, cfg.Session.IdleTimeoutMinutes)
	adminHandler := handler.NewAdminHandler(repos, appService, cfg.Session.IdleTimeoutMinutes)
	healthHandler := handler.NewHealthHandler(db, oidcClient)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS(cfg))

	router.GET("/healthz", healthHandler.Healthz)
	router.GET("/readyz", healthHandler.Readyz)
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.File(cfg.Server.OpenAPIFilePath)
	})

	api := router.Group("/api")
	{
		authGroup := api.Group("/auth")
		authGroup.GET("/login-url", authHandler.LoginURL)
		authGroup.GET("/login", authHandler.Login)
		authGroup.GET("/callback", authHandler.Callback)

		protectedAuth := authGroup.Group("/")
		protectedAuth.Use(middleware.Session(sessionManager))
		protectedAuth.Use(middleware.IdleTimeout(sessionManager, repos.Settings, cfg.Session.IdleTimeoutMinutes))
		protectedAuth.GET("/me", authHandler.Me)
		protectedAuth.POST("/logout", authHandler.Logout)

		protectedPortal := api.Group("/portal")
		protectedPortal.Use(middleware.Session(sessionManager))
		protectedPortal.Use(middleware.IdleTimeout(sessionManager, repos.Settings, cfg.Session.IdleTimeoutMinutes))
		protectedPortal.GET("/apps", portalHandler.Apps)
		protectedPortal.GET("/realms", portalHandler.Realms)
		protectedPortal.GET("/profile", portalHandler.Profile)

		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.Session(sessionManager))
		adminGroup.Use(middleware.IdleTimeout(sessionManager, repos.Settings, cfg.Session.IdleTimeoutMinutes))
		adminGroup.Use(middleware.RequirePortalAdmin())
		adminGroup.GET("/realms", adminHandler.ListRealms)
		adminGroup.GET("/clients", adminHandler.ListClients)
		adminGroup.PUT("/clients/:clientId/meta", adminHandler.UpdateClientMeta)
		adminGroup.GET("/users/:userId", adminHandler.GetUser)
		adminGroup.GET("/settings/session", adminHandler.GetSessionSettings)
		adminGroup.PUT("/settings/session", adminHandler.UpdateSessionSettings)
		adminGroup.GET("/sync-status", adminHandler.SyncStatus)
	}

	server := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("portal-api listening", slog.String("addr", cfg.Server.Addr))
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
