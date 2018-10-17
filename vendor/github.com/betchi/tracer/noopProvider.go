package tracer

import (
	"context"
	"net/http"
)

type noopProvider struct {
	ctx context.Context
}

func (np *noopProvider) NewTracer(config *Config) error {
	return nil
}

func (np *noopProvider) StartTransaction(ctx context.Context, name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{}) {
	return ctx, nil
}

func (np *noopProvider) StartSpan(ctx context.Context, name, spanType string) interface{} {
	return nil
}

func (np *noopProvider) InjectHTTPRequest(span interface{}, req *http.Request) {
}

func (np *noopProvider) SetTag(span interface{}, key string, value interface{}) {
}

func (np *noopProvider) SetHTTPStatusCode(span interface{}, statusCode int) {
}

func (np *noopProvider) SetError(span interface{}, err error) {
}

func (np *noopProvider) Finish(span interface{}) {
}

func (np *noopProvider) CloseTransaction(ctx context.Context) {
}

func (np *noopProvider) Close() {
}
