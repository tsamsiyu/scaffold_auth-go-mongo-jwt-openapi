package tracing

import "context"

const (
	traceIDCtxKey = "traceID"
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDCtxKey, traceID)
}
