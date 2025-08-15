package ctxutils

import (
	"context"
)

type contextKey string

const (
	contextRequestKey contextKey = "requestID"
	contextUserKey    contextKey = "ContextUserKey"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextRequestKey, requestID)
}

func RequestID(ctx context.Context) string {
	return ctx.Value(contextRequestKey).(string)
}

func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, contextUserKey, userID)
}

func UserID(ctx context.Context) uint {
	return ctx.Value(contextUserKey).(uint)
}
