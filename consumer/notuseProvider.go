package consumer

import (
	"context"
)

type notuseProvider struct {
	ctx context.Context
}

func (np notuseProvider) SubscribeMessage() error {
	// Do not process anything
	return nil
}

func (np notuseProvider) UnsubscribeMessage() error {
	// Do not process anything
	return nil
}
