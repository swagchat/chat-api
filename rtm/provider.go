package rtm

import (
	"github.com/swagchat/chat-api/utils"
)

type MessagingInfo struct {
	Message string
}

type provider interface {
	PublishMessage(*MessagingInfo) error
}

func Provider() provider {
	cfg := utils.Config()
	var p provider

	switch cfg.RTM.Provider {
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
