package config

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr string
	LogLevel slog.Level
}

func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		LogLevel: parseLogLevel(getEnv("LOG_LEVEL", "info")),
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
