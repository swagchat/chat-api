package pbroker

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

var KafkaConsumer *kafka.Consumer

type kafkaProvider struct{}

func (kp kafkaProvider) PublishMessage(rtmEvent *RTMEvent) error {
	cfg := utils.Config()
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(rtmEvent)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.PBroker.Kafka.Host, cfg.PBroker.Kafka.Port),
	})
	if err != nil {
		return errors.Wrap(err, "Kafka create producer failure")
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logging.Log(zapcore.ErrorLevel, &logging.AppLog{
						Kind:    "Kafka",
						Message: fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition),
						Error:   err,
					})
				} else {
					logging.Log(zapcore.InfoLevel, &logging.AppLog{
						Kind:    "Kafka",
						Message: fmt.Sprintf("Delivered message to %v\n", ev.TopicPartition),
					})
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := cfg.PBroker.Kafka.Topic
	// for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          buffer.Bytes(),
	}, nil)
	// }

	// Wait for message deliveries
	p.Flush(15 * 1000)

	return nil
}
