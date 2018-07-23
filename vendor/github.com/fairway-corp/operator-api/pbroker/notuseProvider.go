package pbroker

import (
	"context"

	scpb "github.com/swagchat/protobuf"
)

type notuseProvider struct {
	ctx context.Context
}

func (np notuseProvider) PostMessageSwag(m *scpb.Message) error {
	// Do not process anything
	return nil
}

// func (np notuseProvider) PostMessageBot(m *chatpb.BotMessage) error {
// 	// Do not process anything
// 	return nil
// }
