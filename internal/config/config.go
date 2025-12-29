package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cookie   CookieConfig
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig contains database-specific configuration
type DatabaseConfig struct {
	Path            string
	CleanupInterval time.Duration
	RetentionPeriod time.Duration
}

// CookieConfig contains cookie-specific configuration
type CookieConfig struct {
	Prefix   string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite int
}

// New creates a new Config instance with values from environment variables or defaults
func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Path:            getEnv("DB_PATH", "data/exports.db"),
			CleanupInterval: getDurationEnv("DB_CLEANUP_INTERVAL", 24*time.Hour),
			RetentionPeriod: getDurationEnv("DB_RETENTION_PERIOD", 24*time.Hour),
		},
		Cookie: CookieConfig{
			Prefix:   getEnv("COOKIE_PREFIX", "investigator"),
			MaxAge:   getIntEnv("COOKIE_MAX_AGE", 3600*24*30), // 30 days
			HttpOnly: getBoolEnv("COOKIE_HTTP_ONLY", true),
			Secure:   getBoolEnv("COOKIE_SECURE", true),
			SameSite: getIntEnv("COOKIE_SAME_SITE", 3), // http.SameSiteStrictMode = 3
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets an environment variable as int or returns a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getBoolEnv gets an environment variable as bool or returns a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// getDurationEnv gets an environment variable as time.Duration or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}