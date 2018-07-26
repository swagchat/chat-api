package tracer

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/logger"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

type jaegerProvider struct {
	ctx context.Context
}

// Init returns an instance of Jaeger Tracer
func (jp *jaegerProvider) NewTracer(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		logger.Error(err.Error())
		return nil, nil
	}
	return tracer, closer
}
