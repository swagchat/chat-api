package pbroker

import (
	"context"
	"fmt"

	"github.com/fairway-corp/operator-api/utils"
	scpb "github.com/swagchat/protobuf"
)

type provider interface {
	PostMessageSwag(*scpb.Message) error
	// PostMessageBot(*chatpb.BotMessage) error
}

func Provider(ctx context.Context) provider {
	cfg := utils.Config()

	var p provider
	switch cfg.PBroker.Provider {
	case "":
		p = &notuseProvider{
			ctx: ctx,
		}
	case "messageConnector":
		p = &messageConnectorProvider{
			ctx:      ctx,
			endpoint: cfg.PBroker.MessageConnector.Endpoint,
			protocol: cfg.PBroker.MessageConnector.Protocol,
		}
	case "nsq":
		p = &nsqProvider{
			ctx:      ctx,
			endpoint: fmt.Sprintf("%s:%s", cfg.PBroker.NSQ.NsqlookupdHost, cfg.PBroker.NSQ.NsqlookupdPort),
			topic:    cfg.PBroker.NSQ.Topic,
		}
	case "kafka":
		p = &kafkaProvider{
			ctx:      ctx,
			endpoint: fmt.Sprintf("%s:%s", cfg.PBroker.Kafka.Host, cfg.PBroker.Kafka.Port),
			topic:    cfg.PBroker.Kafka.Topic,
		}
	}

	return p
}
