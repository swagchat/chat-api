package sbroker

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

type notuseProvider struct {
	ctx context.Context
}

func (np notuseProvider) SubscribeMessage() error {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "sbroker.notuseProvider.SubscribeMessage")
	defer span.Finish()

	// Do not process anything
	return nil
}

func (np notuseProvider) UnsubscribeMessage() error {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "sbroker.notuseProvider.UnsubscribeMessage")
	defer span.Finish()

	// Do not process anything
	return nil
}
