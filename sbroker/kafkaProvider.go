package sbroker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
	"go.uber.org/zap/zapcore"
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
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:  "messaging-subscribe",
			Error: errors.Wrap(err, "Createt kafka consumer failure"),
		})
		return errors.Wrap(err, "")
	}

	topic := c.SBroker.Kafka.Topic
	err = KafkaConsumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:  "messaging-subscribe",
			Error: errors.Wrap(err, "Subscribe topic failure."),
		})
		return errors.Wrap(err, "")
	}
	_, file, line, _ := runtime.Caller(0)
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:    "messaging-subscribe",
		Message: fmt.Sprintf("%s group.id[%s] topic[%s]", KafkaConsumer.String(), hostname, topic),
		File:    file,
		Line:    line + 2,
	})

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			run = false
			_, file, line, _ := runtime.Caller(0)
			logging.Log(zapcore.InfoLevel, &logging.AppLog{
				Kind:    "messaging-subscribe-terminated",
				Message: fmt.Sprintf("terminated by %s", sig.String()),
				File:    file,
				Line:    line + 1,
			})
		default:
			ev := KafkaConsumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				_, file, line, _ := runtime.Caller(0)
				logging.Log(zapcore.InfoLevel, &logging.AppLog{
					Kind:    "messaging-subscribe-receive",
					Message: fmt.Sprintf("Receive a message"),
					File:    file,
					Line:    line + 1,
				})

				var sm *scpb.Message
				err := json.Unmarshal(e.Value, &sm)
				if err != nil {
					logging.Log(zapcore.ErrorLevel, &logging.AppLog{
						Kind:  "messaging-subscribe",
						Error: errors.Wrap(err, ""),
					})
				}
				log.Printf("%#v\n", sm)
				ctx := context.Background()
				ctx = context.WithValue(ctx, utils.CtxRealm, sm.Workspace)
				ctx = context.WithValue(ctx, utils.CtxUserID, sm.UserId)
				realm := ctx.Value(utils.CtxRealm).(string)
				userID := ctx.Value(utils.CtxUserID).(string)

				var msg *models.Message
				err = json.Unmarshal(e.Value, &msg)
				if err != nil {
					logging.Log(zapcore.ErrorLevel, &logging.AppLog{
						Kind:  "messaging-subscribe",
						Error: errors.Wrap(err, ""),
					})
				}

				msgs := []*models.Message{msg}

				services.PostMessage(ctx, &models.Messages{
					Messages: msgs,
				})
			case kafka.PartitionEOF:
				_, file, line, _ := runtime.Caller(0)
				logging.Log(zapcore.InfoLevel, &logging.AppLog{
					Kind:    "messaging-subscribe",
					Message: e.String(),
					File:    file,
					Line:    line + 1,
				})
			case kafka.Error:
				run = false
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Kind:    "messaging-subscribe",
					Message: e.String(),
				})
				return errors.Wrap(fmt.Errorf("%s", e.String()), "")
			default:
				_, file, line, _ := runtime.Caller(0)
				logging.Log(zapcore.InfoLevel, &logging.AppLog{
					Kind:    "messaging-subscribe",
					Message: e.String(),
					File:    file,
					Line:    line + 1,
				})
			}
		}
	}

	KafkaConsumer.Close()
	_, file, line, _ = runtime.Caller(0)
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:    "messaging-subscribe",
		Message: "close",
		File:    file,
		Line:    line + 1,
	})

	return nil
}

func (kp *kafkaProvider) UnsubscribeMessage() error {
	if KafkaConsumer == nil {
		return nil
	}

	_, file, line, _ := runtime.Caller(0)
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind: "messaging-unsubscribe",
		File: file,
		Line: line + 1,
	})
	return KafkaConsumer.Unsubscribe()
}
