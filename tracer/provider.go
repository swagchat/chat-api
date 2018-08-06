package tracer

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/config"
)

type provider interface {
	NewTracer() error
	StartTransaction(name, transactionType string, opts ...StartTransactionOption) context.Context
	StartSpan(name, spanType string) interface{}
	SetTag(key string, value interface{})
	SetHTTPStatusCode(statusCode int)
	SetUserID(id string)
	Finish(span interface{})
	CloseTransaction()
	Close()
}

type startTransactionOptions struct {
	r *http.Request
}

type StartTransactionOption func(*startTransactionOptions)

func StartTransactionOptionWithHTTPRequest(r *http.Request) StartTransactionOption {
	return func(ops *startTransactionOptions) {
		ops.r = r
	}
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
