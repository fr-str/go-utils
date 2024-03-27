package log

import (
	"log/slog"
	"strings"

	"github.com/fr-str/go-utils/env"
)

// CorrelationIDKey is a default key name for go contexts
const CorrelationIDKey ContextMeta = "correlation_id"
const CorrelationIDHeaderKey = "X-Correlation-ID"

// callerID is a key for context used by TextHandler
const callerID ContextMeta = "callerID"

// Logging configuration parameters.
var (
	logLevel     string = env.Get("LOG_LEVEL", "debug")
	logFormat    string = env.Get("LOG_FORMAT", "json")
	logMultiline bool   = env.Get("LOG_MULTILINE", false)

	// enviroment represents the current environment type.
	enviroment string = env.Get("ENVIRONMENT", "prod") // prod, dev, test
)

// defaultLogger is a default defaultLogger for this log package
var defaultLogger Logger = New()

type Config struct {
	Level     string
	Multiline bool

	// Additional metadata
	Enviroment string
}

// defaults sets default value for fields that have not been set
func (cfg *Config) defaults() {
	if cfg.Level == "" {
		cfg.Level = logLevel
	}
	if cfg.Enviroment == "" {
		cfg.Enviroment = enviroment
	}
}

// slogLevel returns slog.Level
func (cfg Config) slogLevel() slog.Level {
	switch strings.ToLower(cfg.Level)[0] {
	case 't':
		return LevelTrace
	case 'd':
		return slog.LevelDebug
	case 'i':
		return slog.LevelInfo
	case 'w':
		return slog.LevelWarn
	case 'e':
		return slog.LevelError
	case 'f':
		return LevelFatal
	default:
		return slog.LevelInfo
	}
}
