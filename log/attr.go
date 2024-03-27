package log

import (
	"log/slog"
	"time"
)

// An slog.Attr is a key-value pair.

// Any returns any slog.Attr
//
// Use only if value is of a custom type or if specific Attr does not exist
func Any(key string, value any) slog.Attr {
	return slog.Attr{Key: key, Value: slog.AnyValue(value)}
}

// Bool returns bool slog.Attr
func Bool(key string, v bool) slog.Attr {
	return slog.Attr{Key: key, Value: slog.BoolValue(v)}
}

// Duration returns duration slog.Attr
func Duration(key string, v time.Duration) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.StringValue(v.Round(time.Microsecond).String()),
	}
}

// Float64 returns float64 slog.Attr
func Float64(key string, v float64) slog.Attr {
	return slog.Attr{Key: key, Value: slog.Float64Value(v)}
}

// Group returns group slog.Attr
func Group(key string, v ...slog.Attr) slog.Attr {
	return slog.Attr{Key: key, Value: slog.GroupValue(v...)}
}

// Int64 returns int64 slog.Attr
func Int64(key string, value int64) slog.Attr {
	return slog.Attr{Key: key, Value: slog.Int64Value(value)}
}

// Int returns int slog.Attr
func Int(key string, value int) slog.Attr {
	return slog.Attr{Key: key, Value: slog.IntValue(value)}
}

// String returns string slog.Attr
func String(key, value string) slog.Attr {
	return slog.Attr{Key: key, Value: slog.StringValue(value)}
}

// Time returns time slog.Attr
func Time(key string, v time.Time) slog.Attr {
	return slog.Attr{Key: key, Value: slog.TimeValue(v)}
}

// Uint64 returns uint64 slog.Attr
func Uint64(key string, v uint64) slog.Attr {
	return slog.Attr{Key: key, Value: slog.Uint64Value(v)}
}
