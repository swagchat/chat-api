package sbroker

import (
	"context"
	"fmt"
	"net/http"
	"os"

	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/tracer"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct {
	ctx context.Context
}

func (np *nsqProvider) SubscribeMessage() error {
	span := tracer.Provider(np.ctx).StartSpan("SubscribeMessage", "sbroker")
	defer tracer.Provider(np.ctx).Finish(span)

	c := config.Config()
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
			logger.Error(err.Error())
			tracer.Provider(np.ctx).SetError(span, err)
			return err
		}

		logger.Info(fmt.Sprintf("%p", NSQConsumer))

		NSQConsumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			// TODO
			return nil
		}))
		err = NSQConsumer.ConnectToNSQLookupd(c.SBroker.NSQ.NsqlookupdHost + ":" + c.SBroker.NSQ.NsqlookupdPort)
		if err != nil {
			logger.Error(err.Error())
			tracer.Provider(np.ctx).SetError(span, err)
			return err
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

	c := config.Config()
	hostname, err := os.Hostname()
	_, err = http.Post("http://"+c.SBroker.NSQ.NsqdHost+":"+c.SBroker.NSQ.NsqdPort+"/channel/delete?topic="+c.SBroker.NSQ.Topic+"&channel="+hostname, "text/plain", nil)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
