package logger

import (
	"log/slog"
	"time"
)

func String(key, value string) any {
	return slog.String(key, value)
}

func Int(key string, value int) any {
	return slog.Int(key, value)
}

func Int64(key string, value int64) any {
	return slog.Int64(key, value)
}

func Uint64(key string, value uint64) any {
	return slog.Uint64(key, value)
}

func Bool(key string, value bool) any {
	return slog.Bool(key, value)
}

func Time(key string, value time.Time) any {
	return slog.Time(key, value)
}

func Duration(key string, value time.Duration) any {
	return slog.Duration(key, value)
}

func Error(err error) any {
	return String("error", err.Error())
}
