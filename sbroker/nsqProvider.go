package sbroker

import (
	"context"
	"fmt"
	"net/http"
	"os"

	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/tracer"
	"github.com/swagchat/chat-api/utils"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct {
	ctx context.Context
}

func (np *nsqProvider) SubscribeMessage() error {
	span := tracer.Provider(np.ctx).StartSpan("SubscribeMessage", "sbroker")
	defer tracer.Provider(np.ctx).Finish(span)

	c := utils.Config()
	if c.SBroker.NSQ.NsqlookupdHost != "" {
		config := nsq.NewConfig()
		channel := c.SBroker.NSQ.Channel
		hostname, err := os.Hostname()
		if err == nil {
			config.Hostname = hostname
			channel = hostname
		}

		NSQConsumer, err = nsq.NewConsumer(c.SBroker.NSQ.Topic, channel, config)
		if err != nil {
			return errors.Wrap(err, "")
		}

		logger.Info(fmt.Sprintf("%p", NSQConsumer))

		NSQConsumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			// TODO
			return nil
		}))
		err = NSQConsumer.ConnectToNSQLookupd(c.SBroker.NSQ.NsqlookupdHost + ":" + c.SBroker.NSQ.NsqlookupdPort)
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	return nil
}

func (np *nsqProvider) UnsubscribeMessage() error {
	span := tracer.Provider(np.ctx).StartSpan("UnsubscribeMessage", "sbroker")
	defer tracer.Provider(np.ctx).Finish(span)

	if NSQConsumer == nil {
		return nil
	}

	c := utils.Config()
	hostname, err := os.Hostname()
	_, err = http.Post("http://"+c.SBroker.NSQ.NsqdHost+":"+c.SBroker.NSQ.NsqdPort+"/channel/delete?topic="+c.SBroker.NSQ.Topic+"&channel="+hostname, "text/plain", nil)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
