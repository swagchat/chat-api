package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/betchi/zapper"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
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
	span := tracer.Provider(kp.ctx).StartSpan("SubscribeMessage", "consumer")
	defer tracer.Provider(kp.ctx).Finish(span)

	cfg := config.Config()

	host := cfg.Consumer.Kafka.Host
	if host == "" {
		err := errors.New("consumer.kafka.host is empty")
		logger.Error(err.Error())
		tracer.Provider(kp.ctx).SetError(span, err)
		return err
	}

	port := cfg.Consumer.Kafka.Port
	if port == "" {
		err := errors.New("consumer.kafka.port is empty")
		logger.Error(err.Error())
		tracer.Provider(kp.ctx).SetError(span, err)
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
		tracer.Provider(kp.ctx).SetError(span, err)
		return err
	}

	topic := cfg.Consumer.Kafka.Topic
	err = client.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(kp.ctx).SetError(span, err)
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
					tracer.Provider(kp.ctx).SetError(span, err)
					break
				}

				var pbMsg scpb.Message
				err = json.Unmarshal(kafkaMsg.Value, &pbMsg)
				if err != nil {
					logger.Error(err.Error())
					tracer.Provider(kp.ctx).SetError(span, err)
					break
				}
				payload := model.JSONText{}
				err := payload.UnmarshalJSON(pbMsg.Payload)
				if err != nil {
					logger.Error(err.Error())
					tracer.Provider(kp.ctx).SetError(span, err)
					break
				}

				msg := &model.Message{pbMsg, payload}
				req := msg.ConvertToSendMessageRequest()

				ctx := context.Background()
				ctx = context.WithValue(ctx, config.CtxUserID, msg.UserID)

				workspace := cfg.Datastore.Database
				for _, v := range kafkaMsg.Headers {
					logger.Debug(fmt.Sprintf("kafka header %s=%s", v.Key, string(v.Value)))
					if v.Key == config.HeaderWorkspace {
						workspace = string(v.Value)
					}
				}
				ctx = context.WithValue(ctx, config.CtxWorkspace, workspace)

				service.SendMessage(ctx, req)

			case kafka.PartitionEOF:
				logger.Info(e.String())
			case kafka.Error:
				run = false
				err = fmt.Errorf(e.String())
				logger.Error(err.Error())
				tracer.Provider(kp.ctx).SetError(span, err)
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
	span := tracer.Provider(kp.ctx).StartSpan("UnsubscribeMessage", "consumer")
	defer tracer.Provider(kp.ctx).Finish(span)

	if client == nil {
		return nil
	}

	logger.Info("kafka unsubscribe")
	return client.Unsubscribe()
}
