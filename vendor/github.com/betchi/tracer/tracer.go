package tracer

import (
	"context"
	"net/http"

	jaeger "github.com/uber/jaeger-client-go"
)

// Tracer is a interface of tracer
type Tracer interface {
	NewTracer(config *Config) error
	StartTransaction(ctx context.Context, name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{})
	StartSpan(ctx context.Context, name, spanType string) interface{}
	InjectHTTPRequest(span interface{}, req *http.Request)
	SetTag(span interface{}, key string, value interface{})
	SetHTTPStatusCode(span interface{}, statusCode int)
	SetError(span interface{}, err error)
	Finish(span interface{})
	CloseTransaction(ctx context.Context)
	Close()
}

var globalTracer Tracer

type startTransactionOptions struct {
	r *http.Request
}

type StartTransactionOption func(*startTransactionOptions)

func StartTransactionOptionWithHTTPRequest(r *http.Request) StartTransactionOption {
	return func(ops *startTransactionOptions) {
		ops.r = r
	}
}

// InitGlobalTracer initialize global tracer
func InitGlobalTracer(config *Config) error {
	switch config.Provider {
	case "jaeger":
		if config.Jaeger.Logger == nil {
			config.Jaeger.Logger = jaeger.StdLogger
		}
		globalTracer = &jaegerProvider{
			logger: config.Jaeger.Logger,
		}
	case "zipkin":
		if config.Zipkin.Logger == nil {
			config.Zipkin.Logger = jaeger.StdLogger
		}
		globalTracer = &zipkinProvider{
			logger:    config.Zipkin.Logger,
			endpoint:  config.Zipkin.Endpoint,
			batchSize: config.Zipkin.BatchSize,
			timeout:   config.Zipkin.Timeout,
		}
	case "elasticapm":
		globalTracer = &elasticapmProvider{
			logger: config.ElasticAPM.Logger,
		}
	default:
		globalTracer = &noopProvider{}
	}
	return globalTracer.NewTracer(config)
}

// GlobalTracer retrieve global tracer
func GlobalTracer() Tracer {
	return globalTracer
}

// StartTransaction sets a transaction
func StartTransaction(ctx context.Context, name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{}) {
	return globalTracer.StartTransaction(ctx, name, transactionType, opts...)
}

// StartSpan sets a span
func StartSpan(ctx context.Context, name, spanType string) interface{} {
	return globalTracer.StartSpan(ctx, name, spanType)
}

// InjectHTTPRequest injects http request
func InjectHTTPRequest(span interface{}, req *http.Request) {
	globalTracer.InjectHTTPRequest(span, req)
}

// SetTag sets a tag
func SetTag(span interface{}, key string, value interface{}) {
	globalTracer.SetTag(span, key, value)
}

// SetHTTPStatusCode sets a http status code
func SetHTTPStatusCode(span interface{}, statusCode int) {
	globalTracer.SetHTTPStatusCode(span, statusCode)
}

// SetError sets a error
func SetError(span interface{}, err error) {
	globalTracer.SetError(span, err)
}

// Finish makes it finish transaction
func Finish(span interface{}) {
	globalTracer.Finish(span)
}

// CloseTransaction closes transaction
func CloseTransaction(ctx context.Context) {
	globalTracer.CloseTransaction(ctx)
}

// Close closes tracer
func Close() {
	globalTracer.Close()
}
