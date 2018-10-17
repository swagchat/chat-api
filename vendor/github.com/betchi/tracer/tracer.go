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

var tracer Tracer

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
		tracer = &jaegerProvider{
			logger: config.Jaeger.Logger,
		}
	case "zipkin":
		if config.Zipkin.Logger == nil {
			config.Zipkin.Logger = jaeger.StdLogger
		}
		tracer = &zipkinProvider{
			logger:    config.Zipkin.Logger,
			endpoint:  config.Zipkin.Endpoint,
			batchSize: config.Zipkin.BatchSize,
			timeout:   config.Zipkin.Timeout,
		}
	case "elasticapm":
		tracer = &elasticapmProvider{
			logger: config.ElasticAPM.Logger,
		}
	default:
		tracer = &noopProvider{}
	}
	return tracer.NewTracer(config)
}

// GlobalTracer retrieve global tracer
func GlobalTracer() Tracer {
	return tracer
}

// StartTransaction sets a transaction
func StartTransaction(ctx context.Context, name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{}) {
	return tracer.StartTransaction(ctx, name, transactionType, opts...)
}

// StartSpan sets a span
func StartSpan(ctx context.Context, name, spanType string) interface{} {
	return tracer.StartSpan(ctx, name, spanType)
}

// InjectHTTPRequest injects http request
func InjectHTTPRequest(span interface{}, req *http.Request) {
	tracer.InjectHTTPRequest(span, req)
}

// SetTag sets a tag
func SetTag(span interface{}, key string, value interface{}) {
	tracer.SetTag(span, key, value)
}

// SetHTTPStatusCode sets a http status code
func SetHTTPStatusCode(span interface{}, statusCode int) {
	tracer.SetHTTPStatusCode(span, statusCode)
}

// SetError sets a error
func SetError(span interface{}, err error) {
	tracer.SetError(span, err)
}

// Finish makes it finish transaction
func Finish(span interface{}) {
	tracer.Finish(span)
}

// CloseTransaction closes transaction
func CloseTransaction(ctx context.Context) {
	tracer.CloseTransaction(ctx)
}

// Close closes tracer
func Close() {
	tracer.Close()
}
