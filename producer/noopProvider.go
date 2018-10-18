package producer

import (
	"context"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type noopProvider struct {
	ctx context.Context
}

func (np noopProvider) PublishMessage(rtmEvent *scpb.EventData) error {
	// Do not process anything
	return nil
}
