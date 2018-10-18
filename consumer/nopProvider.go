package consumer

import (
	"context"
)

type noopProvider struct {
	ctx context.Context
}

func (np noopProvider) SubscribeMessage() error {
	// Do not process anything
	return nil
}

func (np noopProvider) UnsubscribeMessage() error {
	// Do not process anything
	return nil
}
