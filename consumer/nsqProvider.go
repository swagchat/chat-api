package consumer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	logger "github.com/betchi/zapper"
	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/tracer"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct {
	ctx context.Context
}

func (np *nsqProvider) SubscribeMessage() error {
	span := tracer.Provider(np.ctx).StartSpan("SubscribeMessage", "consumer")
	defer tracer.Provider(np.ctx).Finish(span)

	c := config.Config()
	if c.Consumer.NSQ.NsqlookupdHost != "" {
		config := nsq.NewConfig()
		channel := c.Consumer.NSQ.Channel
		hostname, err := os.Hostname()
		if err == nil {
			config.Hostname = hostname
			channel = hostname
		}

		NSQConsumer, err = nsq.NewConsumer(c.Consumer.NSQ.Topic, channel, config)
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
		err = NSQConsumer.ConnectToNSQLookupd(c.Consumer.NSQ.NsqlookupdHost + ":" + c.Consumer.NSQ.NsqlookupdPort)
		if err != nil {
			logger.Error(err.Error())
			tracer.Provider(np.ctx).SetError(span, err)
			return err
		}
	}

	return nil
}

func (np *nsqProvider) UnsubscribeMessage() error {
	span := tracer.Provider(np.ctx).StartSpan("UnsubscribeMessage", "consumer")
	defer tracer.Provider(np.ctx).Finish(span)

	if NSQConsumer == nil {
		return nil
	}

	c := config.Config()
	hostname, err := os.Hostname()
	_, err = http.Post("http://"+c.Consumer.NSQ.NsqdHost+":"+c.Consumer.NSQ.NsqdPort+"/channel/delete?topic="+c.Consumer.NSQ.Topic+"&channel="+hostname, "text/plain", nil)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
