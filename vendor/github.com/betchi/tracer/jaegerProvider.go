package tracer

import (
	"context"
	"fmt"
	"io"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics"
)

var (
	jaegerTracer opentracing.Tracer
	jaegerCloser io.Closer
)

type jaegerProvider struct {
	logger jaeger.Logger

	// zipkin
	endpoint  string
	batchSize int
	timeout   int
}

func (jp *jaegerProvider) NewTracer(config *Config) error {
	var tracer opentracing.Tracer
	var closer io.Closer

	cfg, err := jaegerConfig.FromEnv()
	if err != nil {
		jp.logger.Error(err.Error())
		return err
	}

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	tracer, closer, err = cfg.New(
		fmt.Sprintf("%s-%s", config.ServiceName, config.ServiceVersion),
		jaegerConfig.Logger(jp.logger),
		jaegerConfig.Metrics(metrics.NullFactory),
		jaegerConfig.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegerConfig.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaegerConfig.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		jp.logger.Error(err.Error())
		return err
	}

	opentracing.SetGlobalTracer(tracer)
	jaegerTracer = tracer
	jaegerCloser = closer

	return nil
}

func (jp *jaegerProvider) StartTransaction(ctx context.Context, name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{}) {
	if jaegerTracer == nil {
		return ctx, nil
	}

	opt := startTransactionOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var span opentracing.Span
	if opt.r == nil {
		span = jaegerTracer.StartSpan(name)
	} else {
		spanCtx, _ := jaegerTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(opt.r.Header))
		span = jaegerTracer.StartSpan(name, ext.RPCServerOption(spanCtx))
	}

	ctxWithSpan := opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctxWithSpan, CtxTracerSpan, span)
	return ctx, span
}

func (jp *jaegerProvider) StartSpan(ctx context.Context, name, spanType string) interface{} {
	if jaegerTracer == nil {
		return nil
	}

	span, _ := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("%s.%s", spanType, name))
	return span
}

func (jp *jaegerProvider) InjectHTTPRequest(span interface{}, req *http.Request) {
	if span == nil {
		return
	}

	s := span.(opentracing.Span)
	ext.SpanKindRPCClient.Set(s)
	ext.HTTPUrl.Set(s, req.RequestURI)
	ext.HTTPMethod.Set(s, req.Method)
	s.Tracer().Inject(
		s.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
}

func (jp *jaegerProvider) SetTag(span interface{}, key string, value interface{}) {
	if span == nil {
		return
	}
	span.(opentracing.Span).SetTag(key, value)
}

func (jp *jaegerProvider) SetHTTPStatusCode(span interface{}, statusCode int) {
	if span == nil {
		return
	}
	span.(opentracing.Span).SetTag("http.status_code", statusCode)
}

func (jp *jaegerProvider) SetError(span interface{}, err error) {
	if span == nil {
		return
	}
	span.(opentracing.Span).SetTag("error", true)
	span.(opentracing.Span).SetTag("message", err.Error())
}

func (jp *jaegerProvider) Finish(span interface{}) {
	if span == nil {
		return
	}
	span.(opentracing.Span).Finish()
}

func (jp *jaegerProvider) CloseTransaction(ctx context.Context) {
	span := ctx.Value(CtxTracerSpan)
	if span == nil {
		return
	}
	span.(opentracing.Span).Finish()
}

func (jp *jaegerProvider) Close() {
	if jaegerCloser == nil {
		return
	}
	jaegerCloser.Close()
}
