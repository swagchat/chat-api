package tracer

import (
	"context"
)

type notuseProvider struct {
	ctx context.Context
}

func (np *notuseProvider) NewTracer() error {
	return nil
}

func (np *notuseProvider) StartTransaction(name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{}) {
	return np.ctx, nil
}

func (np *notuseProvider) StartSpan(name, spanType string) interface{} {
	return nil
}

func (np *notuseProvider) SetTag(span interface{}, key string, value interface{}) {
}

func (np *notuseProvider) SetHTTPStatusCode(span interface{}, statusCode int) {
}

func (np *notuseProvider) SetError(span interface{}, err error) {
}

func (np *notuseProvider) Finish(span interface{}) {
}

func (np *notuseProvider) CloseTransaction() {
}

func (np *notuseProvider) Close() {
}
