package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains all runtime configuration for portal-api.
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Mongo    MongoConfig
	Keycloak KeycloakConfig
	Session  SessionConfig
	Sync     SyncConfig
	CORS     CORSConfig
	Log      LogConfig
}

// AppConfig stores app-level metadata.
type AppConfig struct {
	Name string
	Env  string
}

// ServerConfig stores HTTP server configuration.
type ServerConfig struct {
	Addr            string
	PublicAPIURL    string
	PublicWebURL    string
	OpenAPIFilePath string
}

// MongoConfig stores MongoDB connection parameters.
type MongoConfig struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
}

// KeycloakConfig stores OIDC and Admin API configuration.
type KeycloakConfig struct {
	BaseURL                 string
	PublicURL               string
	Realm                   string
	OIDCClientID            string
	OIDCClientSecret        string
	RedirectURL             string
	PostLogoutRedirectURL   string
	AdminClientID           string
	AdminClientSecret       string
	RequestTimeout          time.Duration
	OIDCScopes              []string
}

// SessionConfig stores portal session configuration.
type SessionConfig struct {
	CookieName             string
	SigningKey             string
	Secure                 bool
	HTTPOnly               bool
	SameSite               string
	IdleTimeoutMinutes     int
	AbsoluteTimeoutMinutes int
	StateCookieName        string
	NonceCookieName        string
	StateCookieMaxAge      time.Duration
}

// SyncConfig stores synchronization behavior.
type SyncConfig struct {
	OnLogin        bool
	TimeoutSeconds int
}

// CORSConfig stores CORS behavior.
type CORSConfig struct {
	AllowedOrigins []string
}

// LogConfig stores logging settings.
type LogConfig struct {
	Level string
}

// MustLoad loads configuration from environment variables.
func MustLoad() Config {
	keycloakBaseURL := strings.TrimRight(getEnv("KEYCLOAK_BASE_URL", "http://localhost:8081"), "/")
	keycloakPublicURL := strings.TrimRight(getEnv("KEYCLOAK_PUBLIC_URL", keycloakBaseURL), "/")
	webPublicURL := strings.TrimRight(getEnv("WEB_PUBLIC_URL", "http://localhost:5173"), "/")
	appPublicURL := strings.TrimRight(getEnv("APP_PUBLIC_URL", webPublicURL+"/api"), "/")
	redirectURL := strings.TrimRight(getEnv("KEYCLOAK_REDIRECT_URL", webPublicURL+"/api/auth/callback"), "/")
	postLogoutRedirectURL := strings.TrimRight(getEnv("KEYCLOAK_POST_LOGOUT_REDIRECT_URL", webPublicURL+"/login"), "/")

	return Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "portal"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Server: ServerConfig{
			Addr:            getEnv("APP_ADDR", ":8080"),
			PublicAPIURL:    appPublicURL,
			PublicWebURL:    webPublicURL,
			OpenAPIFilePath: getEnv("OPENAPI_FILE_PATH", "docs/openapi.yaml"),
		},
		Mongo: MongoConfig{
			URI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database:       getEnv("MONGO_DB", "portal"),
			ConnectTimeout: time.Duration(getEnvInt("MONGO_CONNECT_TIMEOUT_SECONDS", 10)) * time.Second,
		},
		Keycloak: KeycloakConfig{
			BaseURL:               keycloakBaseURL,
			PublicURL:             keycloakPublicURL,
			Realm:                 getEnv("KEYCLOAK_REALM", "portal"),
			OIDCClientID:          getEnv("KEYCLOAK_OIDC_CLIENT_ID", "portal-api"),
			OIDCClientSecret:      getEnv("KEYCLOAK_OIDC_CLIENT_SECRET", "portal-api-secret"),
			RedirectURL:           redirectURL,
			PostLogoutRedirectURL: postLogoutRedirectURL,
			AdminClientID:         getEnv("KEYCLOAK_ADMIN_CLIENT_ID", "portal-sync"),
			AdminClientSecret:     getEnv("KEYCLOAK_ADMIN_CLIENT_SECRET", "portal-sync-secret"),
			RequestTimeout:        time.Duration(getEnvInt("SYNC_TIMEOUT_SECONDS", 10)) * time.Second,
			OIDCScopes:            getEnvSlice("KEYCLOAK_OIDC_SCOPES", []string{"openid", "profile", "email"}),
		},
		Session: SessionConfig{
			CookieName:             getEnv("SESSION_COOKIE_NAME", "portal_session"),
			SigningKey:             getEnv("SESSION_SIGNING_KEY", "portal-signing-key"),
			Secure:                 getEnvBool("SESSION_SECURE", false),
			HTTPOnly:               getEnvBool("SESSION_HTTP_ONLY", true),
			SameSite:               getEnv("SESSION_SAME_SITE", "Lax"),
			IdleTimeoutMinutes:     getEnvInt("SESSION_IDLE_TIMEOUT_MINUTES", 15),
			AbsoluteTimeoutMinutes: getEnvInt("SESSION_ABSOLUTE_TIMEOUT_MINUTES", 480),
			StateCookieName:        getEnv("SESSION_STATE_COOKIE_NAME", "portal_oidc_state"),
			NonceCookieName:        getEnv("SESSION_NONCE_COOKIE_NAME", "portal_oidc_nonce"),
			StateCookieMaxAge:      time.Duration(getEnvInt("SESSION_STATE_COOKIE_MAX_AGE_SECONDS", 300)) * time.Second,
		},
		Sync: SyncConfig{
			OnLogin:        getEnvBool("SYNC_ON_LOGIN", true),
			TimeoutSeconds: getEnvInt("SYNC_TIMEOUT_SECONDS", 10),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"}),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "INFO"),
		},
	}
}

// NewLogger returns a structured slog logger.
func NewLogger(level string) *slog.Logger {
	var slogLevel slog.Level
	switch strings.ToUpper(level) {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	}))
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func getEnvBool(key string, fallback bool) bool {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return value
}

func getEnvSlice(key string, fallback []string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}
