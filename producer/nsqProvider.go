package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"unsafe"

	"github.com/betchi/tracer"
	logger "github.com/betchi/zapper"
	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct {
	ctx context.Context
}

func b2s(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func (np nsqProvider) PublishMessage(rtmEvent *scpb.EventData) error {
	span := tracer.StartSpan(np.ctx, "PublishMessage", "producer")
	defer tracer.Finish(span)

	cfg := config.Config()
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(rtmEvent)

	endpoint := cfg.Producer.NSQ.NsqlookupdHost + ":" + cfg.Producer.NSQ.NsqlookupdPort
	url := fmt.Sprintf("%s/pub?topic=%s", endpoint, cfg.Producer.NSQ.Topic)
	resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		err = errors.Wrap(err, "NSQ post failure")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "NSQ response body read failure")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status code[%d]", resp.StatusCode)
		logger.Error(err.Error())
		tracer.SetError(span, err)
	}

	return nil
}
