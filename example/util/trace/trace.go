package trace

import (
	"context"
	"io"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

var closer io.Closer

func init() {
	// TODO 修改配置
	serviceName := ""
	param := 0.9
	agent := "" // host+":"+port

	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: param, // 采样比例
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: agent,
		},
	}

	tracer, c, err := cfg.NewTracer(
		config.Logger(log.NullLogger),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		panic(err)
	}

	closer = c
	opentracing.SetGlobalTracer(tracer)
}

// GetTraceID
func GetTraceID(ctx context.Context) (traceID string) {
	traceID = "no-trace-id"

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	jctx, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return
	}

	traceID = jctx.TraceID().String()
	return
}

// InjectHeader 注入 open tracing 头信息
func InjectHeader(ctx opentracing.SpanContext, req *http.Request) {
	opentracing.GlobalTracer().Inject(
		ctx,
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	jctx, ok := ctx.(jaeger.SpanContext)
	if !ok {
		return
	}

	// TODO
	req.Header["Bili-Trace-Id"] = req.Header["Uber-Trace-Id"]

	// Envoy 使用 Zipkin 风格头信息
	// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/observability/tracing
	req.Header.Set("x-b3-traceid", jctx.TraceID().String())
	req.Header.Set("x-b3-spanid", jctx.SpanID().String())
	req.Header.Set("x-b3-parentspanid", jctx.ParentID().String())
	if jctx.IsSampled() {
		req.Header.Set("x-b3-sampled", "1")
	}
	if jctx.IsDebug() {
		req.Header.Set("x-b3-flags", "1")
	}
}

// StartFollowSpanFromContext 开启一个follow类型span
// follow类型用于异步任务，可能在root span结束之后才完成
func StartFollowSpanFromContext(ctx context.Context, operation string) (opentracing.Span, context.Context) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return opentracing.StartSpanFromContext(ctx, operation)
	}

	return opentracing.StartSpanFromContext(ctx, operation, opentracing.FollowsFrom(span.Context()))
}

// Stop 停止 trace 协程
func Stop() {
	closer.Close()
}
