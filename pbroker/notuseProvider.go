package pbroker

import (
	"context"
)

type notuseProvider struct {
	ctx context.Context
}

func (np notuseProvider) PublishMessage(rtmEvent *RTMEvent) error {
	// Do not process anything
	return nil
}
