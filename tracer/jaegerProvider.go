package tracer

import (
	"context"
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
	jaeger "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

var (
	jaegerTracer opentracing.Tracer
	jaegerCloser io.Closer
)

type jaegerProvider struct {
	ctx context.Context
}

func (jp *jaegerProvider) NewTracer() error {
	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(fmt.Sprintf("%s:%s", config.AppName, config.BuildVersion), jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	jaegerTracer = tracer
	jaegerCloser = closer
	return nil
}

func (jp *jaegerProvider) StartTransaction(name, transactionType string, opts ...StartTransactionOption) context.Context {
	if jaegerTracer == nil {
		return jp.ctx
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

	ctx := opentracing.ContextWithSpan(jp.ctx, span)
	ctx = context.WithValue(ctx, config.CtxTracerSpan, span)
	return ctx
}

func (jp *jaegerProvider) StartSpan(name, spanType string) interface{} {
	if jaegerTracer == nil {
		return nil
	}

	span, _ := opentracing.StartSpanFromContext(jp.ctx, fmt.Sprintf("%s.%s", spanType, name))
	return span
}

func (jp *jaegerProvider) SetTag(key string, value interface{}) {
	span := jp.ctx.Value(config.CtxTracerSpan)
	if span != nil {
		span.(opentracing.Span).SetTag(key, value)
	}
}

func (jp *jaegerProvider) SetHTTPStatusCode(statusCode int) {
	span := jp.ctx.Value(config.CtxTracerSpan)
	if span != nil {
		span.(opentracing.Span).SetTag("http.status_code", statusCode)
	}
}

func (jp *jaegerProvider) SetUserID(id string) {
	span := jp.ctx.Value(config.CtxTracerSpan)
	if span != nil {
		span.(opentracing.Span).SetTag("app.user_id", id)
	}
}

func (jp *jaegerProvider) Finish(span interface{}) {
	if span != nil {
		span.(opentracing.Span).Finish()
	}
}

func (jp *jaegerProvider) CloseTransaction() {
	span := jp.ctx.Value(config.CtxTracerSpan)
	if span != nil {
		span.(opentracing.Span).Finish()
	}
}

func (jp *jaegerProvider) Close() {
	if jaegerCloser != nil {
		jaegerCloser.Close()
	}
}
