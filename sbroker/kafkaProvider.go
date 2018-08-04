package sbroker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

var client *kafka.Consumer

type kafkaProvider struct {
	ctx context.Context
}

func (kp *kafkaProvider) SubscribeMessage() error {
	span := tracer.Provider(kp.ctx).StartSpan("SubscribeMessage", "sbroker")
	defer tracer.Provider(kp.ctx).Finish(span)

	cfg := utils.Config()

	host := cfg.SBroker.Kafka.Host
	if host == "" {
		err := errors.New("sbroker.kafka.host is empty")
		logger.Error(err.Error())
		return err
	}

	port := cfg.SBroker.Kafka.Port
	if port == "" {
		err := errors.New("sbroker.kafka.port is empty")
		logger.Error(err.Error())
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = utils.GenerateUUID()
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	client, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", host, port),
		"group.id":          hostname,
		// "session.timeout.ms":   6000,
		// "default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}
	})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	topic := cfg.SBroker.Kafka.Topic
	err = client.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Info(fmt.Sprintf("%s group.id[%s] topic[%s]", client.String(), hostname, topic))

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			run = false
			logger.Info(fmt.Sprintf("terminated by %s", sig.String()))
		default:
			ev := client.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				logger.Info("Receive a message")
				kafkaMsg := e
				if err != nil {
					logger.Error(err.Error())
					break
				}

				var pbMsg scpb.Message
				err = json.Unmarshal(kafkaMsg.Value, &pbMsg)
				if err != nil {
					logger.Error(err.Error())
					break
				}
				payload := model.JSONText{}
				err := payload.UnmarshalJSON(pbMsg.Payload)
				if err != nil {
					logger.Error(err.Error())
					break
				}

				msg := &model.Message{pbMsg, payload}
				req := msg.ConvertToCreateMessageRequest()

				ctx := context.Background()
				ctx = context.WithValue(ctx, utils.CtxUserID, msg.UserID)

				workspace := cfg.Datastore.Database
				for _, v := range kafkaMsg.Headers {
					logger.Debug(fmt.Sprintf("kafka header %s=%s", v.Key, string(v.Value)))
					if v.Key == utils.HeaderWorkspace {
						workspace = string(v.Value)
					}
				}
				ctx = context.WithValue(ctx, utils.CtxWorkspace, workspace)

				service.CreateMessage(ctx, req)

			case kafka.PartitionEOF:
				logger.Info(e.String())
			case kafka.Error:
				run = false
				logger.Error(e.String())
			default:
				logger.Info(e.String())
			}
		}
	}

	client.Close()
	logger.Info("kafka close")

	return nil
}

func (kp *kafkaProvider) UnsubscribeMessage() error {
	span := tracer.Provider(kp.ctx).StartSpan("UnsubscribeMessage", "sbroker")
	defer tracer.Provider(kp.ctx).Finish(span)

	if client == nil {
		return nil
	}

	logger.Info("kafka unsubscribe")
	return client.Unsubscribe()
}
