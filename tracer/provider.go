package tracer

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/utils"
)

type provider interface {
	NewTracer(service string) (opentracing.Tracer, io.Closer)
}

func Provider(ctx context.Context) provider {
	cfg := utils.Config()
	var p provider

	switch cfg.Tracer.Provider {
	case "":
		p = &notuseProvider{
			ctx: ctx,
		}
	case "jaeger":
		p = &jaegerProvider{
			ctx: ctx,
		}
	}

	return p
}
