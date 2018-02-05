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
	switch cfg.Rtm.Provider {
	case "":
		provider = &NotUseProvider{}
	case "direct":
		provider = &DirectProvider{}
	case "nsq":
		provider = &NsqProvider{}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "utils.Cfg.Rtm.Provider is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}
