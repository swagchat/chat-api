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
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

var KafkaConsumer *kafka.Consumer

type kafkaProvider struct{}

func (kp *kafkaProvider) SubscribeMessage() error {
	c := utils.Config()

	host := c.SBroker.Kafka.Host
	port := c.SBroker.Kafka.Port

	if host == "" || port == "" {
		return nil
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	var hostname string
	hostname, err := os.Hostname()
	if err != nil {
		hostname = utils.GenerateUUID()
	}
	KafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", host, port),
		"group.id":          hostname,
		// "session.timeout.ms":   6000,
		// "default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}
	})
	if err != nil {
		return errors.Wrap(err, "")
	}

	topic := c.SBroker.Kafka.Topic
	err = KafkaConsumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return errors.Wrap(err, "")
	}
	logger.Info(fmt.Sprintf("%s group.id[%s] topic[%s]", KafkaConsumer.String(), hostname, topic))

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			run = false
			logger.Info(fmt.Sprintf("terminated by %s", sig.String()))
		default:
			ev := KafkaConsumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				logger.Info("Receive a message")
				var sm *scpb.Message
				err := json.Unmarshal(e.Value, &sm)
				if err != nil {
					return errors.Wrap(err, "")
				}

				ctx := context.Background()
				ctx = context.WithValue(ctx, utils.CtxWorkspace, sm.Workspace)
				ctx = context.WithValue(ctx, utils.CtxUserID, sm.UserId)

				var msg *model.Message
				err = json.Unmarshal(e.Value, &msg)
				if err != nil {
					return errors.Wrap(err, "")
				}

				msgs := []*model.Message{msg}

				service.PostMessage(ctx, &model.Messages{
					Messages: msgs,
				})
			case kafka.PartitionEOF:
				logger.Info(e.String())
			case kafka.Error:
				run = false
				return errors.Wrap(fmt.Errorf("%s", e.String()), "")
			default:
				logger.Info(e.String())
			}
		}
	}

	KafkaConsumer.Close()
	logger.Info("close")

	return nil
}

func (kp *kafkaProvider) UnsubscribeMessage() error {
	if KafkaConsumer == nil {
		return nil
	}

	logger.Info("kafka unsubscribe")
	return KafkaConsumer.Unsubscribe()
}
