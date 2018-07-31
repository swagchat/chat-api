package pbroker

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

type notuseProvider struct {
	ctx context.Context
}

func (np notuseProvider) PublishMessage(rtmEvent *RTMEvent) error {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "pbroker.notuseProvider.PublishMessage")
	defer span.Finish()

	// Do not process anything
	return nil
}
