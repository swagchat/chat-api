package rtm

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

type kafkaProvider struct{}

func (kp kafkaProvider) PublishMessage(mi *MessagingInfo) error {
	rawIn := json.RawMessage(mi.Message)
	input, err := rawIn.MarshalJSON()
	cfg := utils.Config()

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.RTM.Kafka.Host, cfg.RTM.Kafka.Port),
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
						Message: fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition),
						Error:   err,
					})
				} else {
					logging.Log(zapcore.ErrorLevel, &logging.AppLog{
						Message: fmt.Sprintf("Delivered message to %v\n", ev.TopicPartition),
						Error:   err,
					})
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := cfg.RTM.Kafka.Topic
	// for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          input,
	}, nil)
	// }

	// Wait for message deliveries
	p.Flush(15 * 1000)
	return nil
}
