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

func GetMessagingProvider() Provider {
	cfg := utils.GetConfig()

	var provider Provider
	switch cfg.RTM.Provider {
	case "":
		provider = &NotUseProvider{}
	case "direct":
		provider = &DirectProvider{}
	case "nsq":
		provider = &NsqProvider{}
	case "kafka":
		provider = &KafkaProvider{}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "RTM Provider is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}
