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

func (np *notuseProvider) StartTransaction(name, transactionType string) context.Context {
	return np.ctx
}

func (np *notuseProvider) StartSpan(name, spanType string) interface{} {
	return nil
}

func (np *notuseProvider) SetTag(key string, value interface{}) {
}

func (np *notuseProvider) SetHTTPStatusCode(statusCode int) {
}

func (np *notuseProvider) SetUserID(id string) {
}

func (np *notuseProvider) Finish(span interface{}) {
}

func (np *notuseProvider) CloseTransaction() {
}

func (np *notuseProvider) Close() {
}
