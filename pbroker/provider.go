package pbroker

import (
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

func Provider() provider {
	cfg := utils.Config()

	var p provider
	switch cfg.PBroker.Provider {
	case "":
		p = &notuseProvider{}
	case "direct":
		p = &directProvider{}
	case "nsq":
		p = &nsqProvider{}
	case "kafka":
		p = &kafkaProvider{}
	}

	return p
}
