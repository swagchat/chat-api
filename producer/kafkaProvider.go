package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/betchi/tracer"
	logger "github.com/betchi/zapper"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

var KafkaConsumer *kafka.Consumer

type kafkaProvider struct {
	ctx context.Context
}

func (kp kafkaProvider) PublishMessage(rtmEvent *scpb.EventData) error {
	span := tracer.StartSpan(kp.ctx, "PublishMessage", "producer")
	defer tracer.Finish(span)

	cfg := config.Config()
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(rtmEvent)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.Producer.Kafka.Host, cfg.Producer.Kafka.Port),
	})
	if err != nil {
		err = errors.Wrap(err, "Kafka create producer failure")
		logger.Error(err.Error())
		tracer.SetError(span, err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					err = fmt.Errorf("Delivery failed: %v", ev.TopicPartition)
					logger.Error(err.Error())
					tracer.SetError(span, err)
				} else {
					logger.Info(fmt.Sprintf("Delivered message to %v\n", ev.TopicPartition))
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := cfg.Producer.Kafka.Topic
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
