package sbroker

import (
	"context"
	"fmt"
	"net/http"
	"os"

	nsq "github.com/nsqio/go-nsq"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct {
	ctx context.Context
}

func (np *nsqProvider) SubscribeMessage() error {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "sbroker.nsqProvider.SubscribeMessage")
	defer span.Finish()

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
	span, _ := opentracing.StartSpanFromContext(np.ctx, "sbroker.nsqProvider.UnsubscribeMessage")
	defer span.Finish()

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
