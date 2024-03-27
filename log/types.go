package log

import (
	"context"
	"log/slog"

	"github.com/fatih/color"
)

// Those are missing levels from slog.Logger
const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

// Colors used by TextHandler renderer
var (
	colorCaller = color.New(90)
	colorTrace  = color.New(color.Bold, color.FgHiWhite)
	colorDebug  = color.New(color.Bold, color.FgMagenta)
	colorInfo   = color.New(color.Bold, color.FgHiBlue)
	colorWarn   = color.New(color.Bold, color.FgYellow)
	colorError  = color.New(color.Bold, color.FgRed)
	colorFatal  = color.New(color.Bold, color.BgHiRed)
)

// Metadata is for easy way to add attrs
//
//	log.Info("Message with metadata", log.Metadata{
//		"query":    util.CleanString(sql),
//		"args":     args,
//		"duration": info.Duration.Round(time.Microsecond).String(),
//	})
type Metadata map[string]any

// ContextMeta is a type to be used with go context
// ie. log.CorrelationIdKey is ContextMeta type
//
//	context.WithValue(ctx, log.CorrelationIdKey, correlationID)
type ContextMeta string

// Logger is a custom slog.Logger implementation
type Logger struct {
	logger *slog.Logger
}

// Logs at level TRACE
func (l *Logger) Trace(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelTrace, msg, args...)
}

// Logs at level DEBUG
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Logs at level INFO
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Logs at level WARN
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Logs at level ERROR
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// Logs at level FATAL
func (l *Logger) Fatal(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelFatal, msg, args...)
	panic(msg)
}

// Logs at level TRACE with given context
func (l *Logger) TraceCtx(ctx context.Context, msg string, args ...any) {
	l.logger.Log(ctx, LevelTrace, msg, args...)
}

// Logs at level DEBUG with given context
func (l *Logger) DebugCtx(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

// Logs at level INFO with given context
func (l *Logger) InfoCtx(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// Logs at level WARN with given context
func (l *Logger) WarnCtx(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

// Logs at level ERROR with given context
func (l *Logger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

// Logs at level FATAL with given context
func (l *Logger) FatalCtx(ctx context.Context, msg string, args ...any) {
	l.logger.Log(ctx, LevelFatal, msg, args...)
	panic(msg)
}
