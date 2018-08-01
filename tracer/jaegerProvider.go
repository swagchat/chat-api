package tracer

import (
	"context"
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

var (
	jaegerTracer opentracing.Tracer
	jaegerCloser io.Closer
)

type jaegerProvider struct {
	ctx context.Context
}

func (jp *jaegerProvider) NewTracer() error {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(fmt.Sprintf("%s:%s", utils.AppName, utils.BuildVersion), config.Logger(jaeger.StdLogger))
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	jaegerTracer = tracer
	jaegerCloser = closer
	return nil
}

func (jp *jaegerProvider) StartTransaction(name, spanType string) context.Context {
	if jaegerTracer == nil {
		return jp.ctx
	}

	span := jaegerTracer.StartSpan(fmt.Sprintf("%s:%s", spanType, name))
	ctx := opentracing.ContextWithSpan(jp.ctx, span)
	ctx = context.WithValue(ctx, utils.CtxTracerSpan, span)
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
	span := jp.ctx.Value(utils.CtxTracerSpan)
	if span != nil {
		span.(opentracing.Span).SetTag(key, value)
	}
}

func (jp *jaegerProvider) Finish(span interface{}) {
	if span != nil {
		span.(opentracing.Span).Finish()
	}
}

func (jp *jaegerProvider) Close() {
	span := jp.ctx.Value(utils.CtxTracerSpan)
	if span != nil {
		span.(opentracing.Span).Finish()
	}

	if jaegerCloser != nil {
		jaegerCloser.Close()
	}
}
