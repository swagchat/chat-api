package pbroker

import (
	"context"

	"github.com/swagchat/chat-api/utils"
)

const (
	MessageEvent = "message"
	UserJoin     = "userJoin"
)

type RTMEvent struct {
	Type    string
	Payload []byte
	UserIDs []string
}

type MessagingInfo struct {
	Message string
}

type provider interface {
	PublishMessage(*RTMEvent) error
}

func Provider(ctx context.Context) provider {
	cfg := utils.Config()

	var p provider
	switch cfg.PBroker.Provider {
	case "":
		p = &notuseProvider{
			ctx: ctx,
		}
	case "direct":
		p = &directProvider{
			ctx: ctx,
		}
	case "nsq":
		p = &nsqProvider{
			ctx: ctx,
		}
	case "kafka":
		p = &kafkaProvider{
			ctx: ctx,
		}
	}

	return p
}
