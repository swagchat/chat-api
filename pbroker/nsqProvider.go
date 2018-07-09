package pbroker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"unsafe"

	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/utils"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct{}

func b2s(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func (np nsqProvider) PublishMessage(rtmEvent *RTMEvent) error {
	cfg := utils.Config()
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(rtmEvent)

	endpoint := cfg.PBroker.NSQ.NsqlookupdHost + ":" + cfg.PBroker.NSQ.NsqlookupdPort
	url := fmt.Sprintf("%s/pub?topic=%s", endpoint, cfg.PBroker.NSQ.Topic)
	resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		return errors.Wrap(err, "NSQ post failure")
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "NSQ response body read failure")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(fmt.Errorf("http status code[%d]", resp.StatusCode), "")
	}

	return nil
}
