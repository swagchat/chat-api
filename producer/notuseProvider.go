package producer

import (
	"context"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type notuseProvider struct {
	ctx context.Context
}

func (np notuseProvider) PublishMessage(rtmEvent *scpb.EventData) error {
	// Do not process anything
	return nil
}
