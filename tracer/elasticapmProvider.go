package tracer

import (
	"context"

	elasticapm "github.com/elastic/apm-agent-go"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
)

var elasticapmTracer *elasticapm.Tracer

type elasticapmProvider struct {
	ctx context.Context
}

func (ep *elasticapmProvider) NewTracer() error {
	tracer, err := elasticapm.NewTracer(config.AppName, config.BuildVersion)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	tracer.SetLogger(logger.Logger())
	tracer.SetCaptureBody(elasticapm.CaptureBodyAll)
	elasticapmTracer = tracer
	return nil
}

func (ep *elasticapmProvider) StartTransaction(name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{}) {
	if elasticapmTracer == nil {
		return ep.ctx, nil
	}

	transaction := elasticapmTracer.StartTransaction(name, transactionType)
	ctx := elasticapm.ContextWithTransaction(ep.ctx, transaction)
	ctx = context.WithValue(ctx, config.CtxTracerTransaction, transaction)
	return ctx, transaction
}

func (ep *elasticapmProvider) StartSpan(name, spanType string) interface{} {
	span, _ := elasticapm.StartSpan(ep.ctx, name, spanType)
	return span
}

func (ep *elasticapmProvider) SetTag(span interface{}, key string, value interface{}) {
	transaction := ep.ctx.Value(config.CtxTracerTransaction)
	if transaction != nil {
		txCtx := transaction.(*elasticapm.Transaction).Context
		// txCtx.SetTag(key, fmt.Sprintf("%v", value))
		txCtx.SetCustom(key, value)
	}
}

func (ep *elasticapmProvider) SetHTTPStatusCode(span interface{}, statusCode int) {
	transaction := ep.ctx.Value(config.CtxTracerTransaction)
	if transaction != nil {
		txCtx := transaction.(*elasticapm.Transaction).Context
		txCtx.SetHTTPStatusCode(statusCode)
	}
}

func (ep *elasticapmProvider) SetError(span interface{}, err error) {
	// TODO
}

func (ep *elasticapmProvider) Finish(span interface{}) {
	if span != nil {
		span.(*elasticapm.Span).End()
	}
}

func (ep *elasticapmProvider) CloseTransaction() {
	transaction := ep.ctx.Value(config.CtxTracerTransaction)
	if transaction != nil {
		transaction.(*elasticapm.Transaction).End()
	}

	elasticapmTracer.Flush(nil)
}

func (ep *elasticapmProvider) Close() {
}
