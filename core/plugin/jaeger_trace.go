package plugin

import (
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

// GetTraceIDFromMd 从 metadata map 中提取 trace_id
// metadata 中 uber-trace-id 格式: traceID:spanID:parentID:sampled
func GetTraceIDFromMd(md map[string]string) string {
	if md == nil {
		return ""
	}
	v, ok := md[JAEGER_KEY]
	if !ok || v == "" {
		return ""
	}
	parts := strings.Split(v, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// GetTraceIDFromSpan 从 opentracing.Span 中提取 trace_id
func GetTraceIDFromSpan(span opentracing.Span) string {
	if span == nil {
		return ""
	}
	ctx, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return ""
	}
	return ctx.TraceID().String()
}
