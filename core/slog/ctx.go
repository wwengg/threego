package slog

import "context"

type contextKey string

const traceIDKey contextKey = "trace_id"

// WithTraceID 将 trace_id 存入 context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID 从 context 中取出 trace_id
func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey).(string); ok {
		return v
	}
	return ""
}
