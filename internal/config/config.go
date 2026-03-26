package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server     ServerConfig
	Mongo      MongoConfig
	Keycloak   KeycloakConfig
	Session    SessionConfig
	Permission PermissionConfig
	Log        LogConfig
}

type ServerConfig struct {
	ListenAddr      string
	PublicAPIURL    string
	PublicWebURL    string
	CookieName      string
	CookieSecure    bool
	CookieDomain    string
	CookieSameSite  string
	AllowedOrigins  []string
	OpenAPIFilePath string
}

type MongoConfig struct {
	URI            string
	Database       string
	ConnectTimeout int
}

type KeycloakConfig struct {
	BaseURL            string
	Realm              string
	AdminRealm         string
	ClientID           string
	ClientSecret       string
	AdminClientID      string
	AdminClientSecret  string
	AdminUsername      string
	AdminPassword      string
	RedirectURL        string
	LogoutRedirectURL  string
	Scopes             []string
	RequestTimeoutSecs int
	SkipTLSVerify      bool
}

type SessionConfig struct {
	DefaultIdleTimeoutMinutes int
	AbsoluteTTLHours          int
	StateCookieName           string
	NonceCookieName           string
	StateCookieMaxAgeSeconds  int
}

type PermissionConfig struct {
	AdminRealmRoles []string
}

type LogConfig struct {
	Level string
}

func MustLoad() Config {
	return Config{
		Server: ServerConfig{
			ListenAddr:      getEnv("PORTAL_API_LISTEN_ADDR", ":8080"),
			PublicAPIURL:    getEnv("PORTAL_PUBLIC_API_URL", "http://localhost:8080"),
			PublicWebURL:    getEnv("PORTAL_PUBLIC_WEB_URL", "http://localhost:5173"),
			CookieName:      getEnv("PORTAL_SESSION_COOKIE_NAME", "portal_session"),
			CookieSecure:    getEnvBool("PORTAL_COOKIE_SECURE", false),
			CookieDomain:    getEnv("PORTAL_COOKIE_DOMAIN", ""),
			CookieSameSite:  getEnv("PORTAL_COOKIE_SAMESITE", "Lax"),
			AllowedOrigins:  getEnvSlice("PORTAL_ALLOWED_ORIGINS", []string{"http://localhost:5173"}),
			OpenAPIFilePath: getEnv("PORTAL_OPENAPI_FILE_PATH", "docs/openapi.yaml"),
		},
		Mongo: MongoConfig{
			URI:            getEnv("PORTAL_MONGO_URI", "mongodb://localhost:27017"),
			Database:       getEnv("PORTAL_MONGO_DATABASE", "portal"),
			ConnectTimeout: getEnvInt("PORTAL_MONGO_CONNECT_TIMEOUT_SECS", 10),
		},
		Keycloak: KeycloakConfig{
			BaseURL:            strings.TrimRight(getEnv("KEYCLOAK_BASE_URL", "http://localhost:8081"), "/"),
			Realm:              getEnv("KEYCLOAK_REALM", "portal"),
			AdminRealm:         getEnv("KEYCLOAK_ADMIN_REALM", "master"),
			ClientID:           getEnv("KEYCLOAK_CLIENT_ID", "portal-api"),
			ClientSecret:       getEnv("KEYCLOAK_CLIENT_SECRET", "change-me"),
			AdminClientID:      getEnv("KEYCLOAK_ADMIN_CLIENT_ID", "portal-admin"),
			AdminClientSecret:  getEnv("KEYCLOAK_ADMIN_CLIENT_SECRET", "change-me"),
			AdminUsername:      getEnv("KEYCLOAK_ADMIN_USERNAME", ""),
			AdminPassword:      getEnv("KEYCLOAK_ADMIN_PASSWORD", ""),
			RedirectURL:        getEnv("KEYCLOAK_REDIRECT_URL", "http://localhost:8080/api/v1/auth/callback"),
			LogoutRedirectURL:  getEnv("KEYCLOAK_LOGOUT_REDIRECT_URL", "http://localhost:5173/login"),
			Scopes:             getEnvSlice("KEYCLOAK_SCOPES", []string{"openid", "profile", "email"}),
			RequestTimeoutSecs: getEnvInt("KEYCLOAK_REQUEST_TIMEOUT_SECS", 10),
			SkipTLSVerify:      getEnvBool("KEYCLOAK_SKIP_TLS_VERIFY", false),
		},
		Session: SessionConfig{
			DefaultIdleTimeoutMinutes: getEnvInt("PORTAL_IDLE_TIMEOUT_MINUTES", 15),
			AbsoluteTTLHours:          getEnvInt("PORTAL_SESSION_TTL_HOURS", 8),
			StateCookieName:           getEnv("PORTAL_STATE_COOKIE_NAME", "portal_oidc_state"),
			NonceCookieName:           getEnv("PORTAL_NONCE_COOKIE_NAME", "portal_oidc_nonce"),
			StateCookieMaxAgeSeconds:  getEnvInt("PORTAL_STATE_COOKIE_MAX_AGE_SECONDS", 300),
		},
		Permission: PermissionConfig{
			AdminRealmRoles: getEnvSlice("PORTAL_ADMIN_REALM_ROLES", []string{"portal-admin", "realm-admin"}),
		},
		Log: LogConfig{
			Level: getEnv("PORTAL_LOG_LEVEL", "INFO"),
		},
	}
}

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
