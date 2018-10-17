package tracer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	transportZipkin "github.com/uber/jaeger-client-go/transport/zipkin"
	"github.com/uber/jaeger-client-go/zipkin"
)

var (
	zipkinTracer opentracing.Tracer
	zipkinCloser io.Closer
)

type zipkinProvider struct {
	logger    jaeger.Logger
	endpoint  string
	batchSize int
	timeout   int
}

func (jp *zipkinProvider) NewTracer(config *Config) error {
	var tracer opentracing.Tracer
	var closer io.Closer

	transport, err := transportZipkin.NewHTTPTransport(
		jp.endpoint,
		transportZipkin.HTTPBatchSize(jp.batchSize),
		transportZipkin.HTTPTimeout(time.Second*time.Duration(jp.timeout)),
		transportZipkin.HTTPLogger(jp.logger),
	)
	if err != nil {
		jp.logger.Error(err.Error())
		return err
	}

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	tracer, closer = jaeger.NewTracer(
		fmt.Sprintf("%s-%s", config.ServiceName, config.ServiceVersion),
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
		jaeger.TracerOptions.Gen128Bit(true),
		jaeger.TracerOptions.Logger(jp.logger),
	)

	opentracing.SetGlobalTracer(tracer)
	zipkinTracer = tracer
	zipkinCloser = closer

	return nil
}

func (jp *zipkinProvider) StartTransaction(ctx context.Context, name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{}) {
	if zipkinTracer == nil {
		return ctx, nil
	}

	opt := startTransactionOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var span opentracing.Span
	if opt.r == nil {
		span = zipkinTracer.StartSpan(name)
	} else {
		spanCtx, _ := zipkinTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(opt.r.Header))
		span = zipkinTracer.StartSpan(name, ext.RPCServerOption(spanCtx))
	}

	ctxWithSpan := opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctxWithSpan, CtxTracerSpan, span)
	return ctx, span
}

func (jp *zipkinProvider) StartSpan(ctx context.Context, name, spanType string) interface{} {
	if zipkinTracer == nil {
		return nil
	}

	span, _ := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("%s.%s", spanType, name))
	return span
}

func (jp *zipkinProvider) InjectHTTPRequest(span interface{}, req *http.Request) {
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

func (jp *zipkinProvider) SetTag(span interface{}, key string, value interface{}) {
	if span == nil {
		return
	}
	span.(opentracing.Span).SetTag(key, value)
}

func (jp *zipkinProvider) SetHTTPStatusCode(span interface{}, statusCode int) {
	if span == nil {
		return
	}
	span.(opentracing.Span).SetTag("http.status_code", statusCode)
}

func (jp *zipkinProvider) SetError(span interface{}, err error) {
	if span == nil {
		return
	}
	span.(opentracing.Span).SetTag("error", true)
	span.(opentracing.Span).SetTag("message", err.Error())
}

func (jp *zipkinProvider) Finish(span interface{}) {
	if span == nil {
		return
	}
	span.(opentracing.Span).Finish()
}

func (jp *zipkinProvider) CloseTransaction(ctx context.Context) {
	span := ctx.Value(CtxTracerSpan)
	if span == nil {
		return
	}
	span.(opentracing.Span).Finish()
}

func (jp *zipkinProvider) Close() {
	if zipkinCloser == nil {
		return
	}
	zipkinCloser.Close()
}
