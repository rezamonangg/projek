package common

import (
	"context"
)

type ctxKey int

const logContextKey ctxKey = 0

type LogContext map[string]any

func WithLogContext(ctx context.Context, key string, value any) context.Context {
	existing, _ := ctx.Value(logContextKey).(LogContext)
	if existing == nil {
		existing = LogContext{}
	}
	newCtx := make(LogContext)
	for k, v := range existing {
		newCtx[k] = v
	}
	newCtx[key] = value
	return context.WithValue(ctx, logContextKey, newCtx)
}

func GetLogContext(ctx context.Context) LogContext {
	if ctx == nil {
		return LogContext{}
	}
	if lc, ok := ctx.Value(logContextKey).(LogContext); ok {
		return lc
	}
	return LogContext{}
}
