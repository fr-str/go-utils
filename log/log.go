package log

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// New returns logger with default configuration
func New() Logger {
	return NewWithConfiguration(Config{
		Level:      logLevel,
		Multiline:  logMultiline,
		Enviroment: enviroment,
	})
}

// NewWithConfiguration returns logger with given configuration
func NewWithConfiguration(cfg Config) Logger {
	cfg.defaults()

	var handler slog.Handler

	if strings.ToLower(cfg.Enviroment) == "prod" || logFormat == "json" {
		// Production/JSON logger
		opts := &slog.HandlerOptions{
			Level: cfg.slogLevel(),
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					switch level {
					case LevelTrace:
						a.Value = slog.StringValue("TRACE")
					case LevelFatal:
						a.Value = slog.StringValue("FATAL")
					}
				}
				return a
			},
		}
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Default
		handler = TextHandler(os.Stdout, cfg)
	}

	handler = MetadataHandler(handler, cfg)
	return Logger{
		logger: slog.New(handler),
	}
}

// Logs with default logger at level TRACE
func Trace(msg string, args ...any) {
	defaultLogger.Trace(msg, args...)
}

// Logs with default logger at level DEBUG
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Logs with default logger at level INFO
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Logs with default logger at level WARN
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Logs with default logger at level ERROR
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// Logs with default logger at level FATAL
func Fatal(msg string, args ...any) {
	defaultLogger.Fatal(msg, args...)
}

// Logs with default logger at level TRACE with given context
func TraceCtx(ctx context.Context, msg string, args ...any) {
	defaultLogger.TraceCtx(ctx, msg, args...)
}

// Logs with default logger at level DEBUG with given context
func DebugCtx(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugCtx(ctx, msg, args...)
}

// Logs with default logger at level INFO with given context
func InfoCtx(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoCtx(ctx, msg, args...)
}

// Logs with default logger at level WARN with given context
func WarnCtx(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnCtx(ctx, msg, args...)
}

// Logs with default logger at level ERROR with given context
func ErrorCtx(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorCtx(ctx, msg, args...)
}

// Logs with default logger at level FATAL with given context
func FatalCtx(ctx context.Context, msg string, args ...any) {
	defaultLogger.FatalCtx(ctx, msg, args...)
}
