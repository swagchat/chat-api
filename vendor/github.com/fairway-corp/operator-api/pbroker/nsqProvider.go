package pbroker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"unsafe"

	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	scpb "github.com/swagchat/protobuf"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct {
	ctx      context.Context
	endpoint string
	topic    string
}

func b2s(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func (np nsqProvider) PostMessageSwag(m *scpb.Message) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	return np.send(b)
}

// func (np nsqProvider) PostMessageBot(m *chatpb.BotMessage) error {
// 	b := new(bytes.Buffer)
// 	json.NewEncoder(b).Encode(m)
// 	return np.send(b)
// }

func (np nsqProvider) send(b *bytes.Buffer) error {
	url := fmt.Sprintf("%s/pub?topic=%s", np.endpoint, np.topic)
	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return errors.Wrap(err, "NSQ post failure")
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "NSQ response body read failure")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code[%d]", resp.StatusCode)
	}

	return nil
}
