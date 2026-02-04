package plugin

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/smallnest/rpcx/share"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func NewTracer(servicename string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		/*
		   jaeger.SamplerTypeConst:          全量采集，采样率设置0,1 分别对应打开和关闭
		   jaeger.SamplerTypeProbabilistic:  概率采集，默认万份之一，0~1之间取值，
		   jaeger.SamplerTypeRateLimiting:   限速采集，每秒只能采集一定量的数据
		   jaeger.SamplerTypeRemote:         一种动态采集策略，根据当前系统的访问量调节采集策略
		*/
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	return tracer, closer, err
}

const JAEGER_KEY = "uber-trace-id"

// 只适用于 jaeger
// 原理用 context 传递 "__req_metadata":""uber-trace-id -> 6f8b8a1101b0124f:6f8b8a1101b0124f:0000000000000000:1
//
//	uber-trace-id traceID : spanID : parentID: sampled bool
func GenerateSpanWithContext(ctx context.Context, operationName string) (opentracing.Span, context.Context, error) {
	md := ctx.Value(share.ReqMetaDataKey) // share.ReqMetaDataKey 固定值 "__req_metadata"  可自定义
	var span opentracing.Span

	tracer := opentracing.GlobalTracer()

	if md != nil {
		carrier := opentracing.TextMapCarrier(md.(map[string]string))
		spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			log.Printf("metadata error %s\n", err)
			return nil, nil, err
		}
		span = tracer.StartSpan(operationName, ext.RPCServerOption(spanContext))
	} else {
		span = opentracing.StartSpan(operationName)
	}

	metadata := opentracing.TextMapCarrier(make(map[string]string))
	err := tracer.Inject(span.Context(), opentracing.TextMap, metadata)
	if err != nil {
		return nil, nil, err
	}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, (map[string]string)(metadata))
	return span, ctx, nil
}

// 不破坏原数据，返回追加了Jaeger_key的数据
func GenerateSpanByMd(md map[string]string, operationName string) (opentracing.Span, map[string]string, error) {
	if len(md) == 0 {
		md = make(map[string]string)
	}
	v, ok := md[JAEGER_KEY]
	md2 := make(map[string]string)
	if ok {
		md2[JAEGER_KEY] = v
	}
	if span, md3, err := GenerateSpanWithMap(md2, operationName); err == nil {
		md[JAEGER_KEY] = md3[JAEGER_KEY]
		return span, md, nil
	} else {
		return nil, nil, err
	}
}

func GenerateSpanWithMap(md2 map[string]string, operationName string) (opentracing.Span, map[string]string, error) {
	var span opentracing.Span

	tracer := opentracing.GlobalTracer()

	if md2 != nil {
		carrier := opentracing.TextMapCarrier(md2)
		spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			log.Printf("metadata error %s\n", err)
			return nil, nil, err
		}
		span = tracer.StartSpan(operationName, ext.RPCServerOption(spanContext))
	} else {
		span = opentracing.StartSpan(operationName)
	}

	metadata := opentracing.TextMapCarrier(make(map[string]string))
	err := tracer.Inject(span.Context(), opentracing.TextMap, metadata)
	if err != nil {
		return nil, nil, err
	}

	return span, metadata, nil
}
