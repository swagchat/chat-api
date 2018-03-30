package rtm

import (
	"os"

	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap"
)

type MessagingInfo struct {
	Message string
}

type Provider interface {
	Init() error
	PublishMessage(*MessagingInfo) error
}

func RTMProvider() Provider {
	cfg := utils.Config()

	var p Provider
	switch cfg.RTM.Provider {
	case "":
		p = &NotUseProvider{}
	case "direct":
		p = &DirectProvider{}
	case "nsq":
		p = &NsqProvider{}
	case "kafka":
		p = &KafkaProvider{}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "RTM Provider is incorrect"),
		)
		os.Exit(0)
	}
	return p
}
