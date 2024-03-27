package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

// badKey copied from slog
const badKey = "!BADKEY"

// pathPrefix is a project directory
var pathPrefix string

func init() {
	_, b, _, _ := runtime.Caller(0)
	pathPrefix = filepath.Join(filepath.Dir(b), "../..") + "/"
}

// metadataHandler handles Metadata type and adds attrs to logs
type metadataHandler struct {
	slog.Handler
	cfg Config
}

// MetadataHandler returns slog.Handler with metadata
func MetadataHandler(handler slog.Handler, cfg Config) *metadataHandler {
	return &metadataHandler{
		Handler: handler,
		cfg:     cfg,
	}
}

func (h *metadataHandler) Handle(ctx context.Context, r slog.Record) error {
	caller := h.getCaller()
	attrs := []slog.Attr{
		slog.String("enviroment", h.cfg.Enviroment),
		slog.String("caller", caller),
		slog.String(string(CorrelationIDKey), h.getCorrelationId(ctx)),
	}

	r.Attrs(func(a slog.Attr) bool {
		if a.Key == badKey {
			// Flatten metadata if key is not given
			meta, isMetadata := a.Value.Any().(Metadata)
			if isMetadata {
				for k, v := range meta {
					attrs = append(attrs, slog.Attr{
						Key: k, Value: slog.AnyValue(v),
					})
				}
				return true
			}
		}

		attrs = append(attrs, a)
		return true
	})

	// Recreate record with new attrs
	record := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	record.AddAttrs(attrs...)

	// Add callerID for TextHandler
	ctx = context.WithValue(ctx, callerID, caller)

	return h.Handler.Handle(ctx, record)
}

// getCorrelationId extracts Correlation ID from ctx
func (h *metadataHandler) getCorrelationId(ctx context.Context) string {
	correlationId := "unknown"
	if ctx == nil {
		return correlationId
	}

	correlationId, ok := ctx.Value(CorrelationIDKey).(string)
	if !ok {
		return correlationId
	}

	return correlationId
}

// getCaller returns caller with package, file and line number
func (*metadataHandler) getCaller() string {
	_, path, line, _ := runtime.Caller(6)
	return fmt.Sprintf("%s:%d", strings.TrimPrefix(path, pathPrefix), line)
}

func (*metadataHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (*metadataHandler) WithAttrs([]slog.Attr) slog.Handler {
	return nil
}

func (*metadataHandler) WithGroup(string) slog.Handler {
	return nil
}

// textHandler is a custom slog.Handler which renders logs with correct format
type textHandler struct {
	out io.Writer
	cfg Config

	marshaler func(any) ([]byte, error)
}

// TextHandler returns slog.Handler for rendering message in text format
func TextHandler(out io.Writer, cfg Config) *textHandler {
	// Choose if json.marshal must be idented or not
	marshaler := json.Marshal
	if cfg.Multiline {
		marshaler = func(x any) ([]byte, error) {
			return json.MarshalIndent(x, "", "  ")
		}
	}

	return &textHandler{
		out:       out,
		cfg:       cfg,
		marshaler: marshaler,
	}
}

// Handle renders log message and writes on output
func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	var level string

	switch r.Level {
	case LevelTrace:
		level = colorTrace.Sprint("TRACE")
	case slog.LevelDebug:
		level = colorDebug.Sprint("DEBUG")
	case slog.LevelInfo:
		level = colorInfo.Sprint("INFO ")
	case slog.LevelWarn:
		level = colorWarn.Sprint("WARN ")
	case slog.LevelError:
		level = colorError.Sprint("ERROR")
	case LevelFatal:
		level = colorFatal.Sprint("FATAL")
	}

	// fields Slice -> Map
	var fields = make(map[string]any, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	b, err := h.marshaler(fields)
	if err != nil {
		return err
	}
	if len(b) <= 2 {
		b = nil
	}

	_, err = fmt.Fprintln(h.out,
		r.Time.Format("15:05:05.000"),
		level,
		colorCaller.Sprint(h.getCaller(ctx)),
		color.WhiteString(r.Message),
		color.WhiteString(string(b)),
	)
	return err
}

func (*textHandler) getCaller(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	caller, ok := ctx.Value(callerID).(string)
	if !ok {
		return ""
	}

	before, caller, ok := strings.Cut(caller, "/")
	if !ok {
		caller = before
	}
	return caller
}

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level.Level() >= h.cfg.slogLevel()
}

func (*textHandler) WithAttrs([]slog.Attr) slog.Handler {
	return nil
}

func (*textHandler) WithGroup(string) slog.Handler {
	return nil
}

// httpHandle logging middleware
type httpHandle struct {
	handler http.Handler
}

// HTTPHandler is a logging middleware for http server
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/v1/hello", HelloHandler)
//	log.Fatal(http.ListenAndServe(addr, log.HTTPHandler(mux)))
func HTTPHandler(handler http.Handler) http.Handler {
	return &httpHandle{handler: handler}
}

// ServeHTTP is a http.Handler implementation
func (h *httpHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	correlationID := r.Header.Get(CorrelationIDHeaderKey)
	if correlationID == "" {
		correlationID = fmt.Sprintf("unknown-%s", uuid.NewString())
	}
	ctx := context.WithValue(r.Context(), CorrelationIDKey, correlationID)
	r = r.WithContext(ctx)

	remoteAddr := r.RemoteAddr
	fwdAddr := r.Header.Get("X-Forwarded-For")
	if fwdAddr != "" {
		// Got X-Forwarded-For
		remoteAddr = fwdAddr

		// If we got an array, grab the first IP
		ips := strings.Split(fwdAddr, ", ")
		if len(ips) > 1 {
			remoteAddr = ips[0]
		}
	}

	ts := time.Now()
	h.handler.ServeHTTP(w, r)
	InfoCtx(r.Context(), r.Method,
		String("path", r.URL.Path),
		Duration("duration", time.Since(ts)),
		String("remote", remoteAddr),
		String("agent", r.UserAgent()),
	)
}

func (h *httpHandle) Rank() int8 {
	return 127
}
