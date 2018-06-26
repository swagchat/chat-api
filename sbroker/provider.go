package sbroker

import (
	"github.com/swagchat/chat-api/utils"
)

type provider interface {
	SubscribeMessage() error
	UnsubscribeMessage() error
}

func Provider() provider {
	cfg := utils.Config()

	var p provider
	switch cfg.SBroker.Provider {
	case "":
		p = &notuseProvider{}
	case "nsq":
		p = &nsqProvider{}
	case "kafka":
		p = &kafkaProvider{}
	}

	return p
}
