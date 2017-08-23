package rtm

import (
	"os"

	"github.com/fairway-corp/swagchat-api/utils"
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
	var provider Provider
	switch utils.Cfg.Rtm.Provider {
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
