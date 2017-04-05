package messaging

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
	switch utils.Cfg.ApiServer.Messaging {
	case "":
		provider = &NotUseProvider{}
	case "gcpPubSub":
		provider = &GcpPubSubProvider{
			thumbnailTopic:    utils.Cfg.GcpPubsub.ThumbnailTopic,
			pushMessageTopic:  utils.Cfg.GcpPubsub.PushMessageTopic,
			scope:             utils.Cfg.GcpPubsub.Scope,
			jwtConfigFilepath: utils.Cfg.GcpPubsub.JwtConfigFilepath,
		}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "utils.Cfg.ApiServer.Messaging is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}
