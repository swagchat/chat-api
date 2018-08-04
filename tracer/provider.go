package tracer

import (
	"context"

	"github.com/swagchat/chat-api/config"
)

type provider interface {
	NewTracer() error
	StartTransaction(name, transactionType string) context.Context
	StartSpan(name, spanType string) interface{}
	SetTag(key string, value interface{})
	SetHTTPStatusCode(statusCode int)
	SetUserID(id string)
	Finish(span interface{})
	CloseTransaction()
	Close()
}

func Provider(ctx context.Context) provider {
	cfg := config.Config()
	var p provider

	switch cfg.Tracer.Provider {
	case "jaeger":
		p = &jaegerProvider{
			ctx: ctx,
		}
	case "elasticapm":
		p = &elasticapmProvider{
			ctx: ctx,
		}
	default:
		p = &notuseProvider{
			ctx: ctx,
		}
	}

	return p
}
