package pbroker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/utils"
	scpb "github.com/swagchat/protobuf"
)

var KafkaConsumer *kafka.Consumer

type kafkaProvider struct {
	ctx      context.Context
	protocol string
	endpoint string
	topic    string
}

func (kp kafkaProvider) PostMessageSwag(m *scpb.Message) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	return kp.send(b)
}

// func (kp kafkaProvider) PostMessageBot(m *chatpb.BotMessage) error {
// 	b := new(bytes.Buffer)
// 	json.NewEncoder(b).Encode(m)
// 	return kp.send(b)
// }

func (kp kafkaProvider) send(b *bytes.Buffer) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kp.endpoint,
	})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Info(fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition))
				} else {
					logger.Info(fmt.Sprintf("Delivered message to %v\n", ev.TopicPartition))
				}
			default:
				logger.Info(e.String())
			}
		}
	}()

	topic := kp.topic
	workspaceHeader := kafka.Header{
		Key:   utils.HeaderWorkspace,
		Value: []byte(kp.ctx.Value(utils.CtxWorkspace).(string)),
	}
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b.Bytes(),
		Headers:        []kafka.Header{workspaceHeader},
	}, nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Wait for message deliveries
	p.Flush(15 * 1000)

	return nil
}
