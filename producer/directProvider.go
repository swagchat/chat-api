package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/betchi/tracer"
	logger "github.com/betchi/zapper"
	"github.com/swagchat/chat-api/config"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type directProvider struct {
	ctx context.Context
}

func (dp directProvider) PublishMessage(rtmEvent *scpb.EventData) error {
	span := tracer.StartSpan(dp.ctx, "PublishMessage", "producer")
	defer tracer.Finish(span)

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(rtmEvent)

	endpoint := fmt.Sprintf("%s/message", config.Config().Producer.Direct.Endpoint)

	req, err := http.NewRequest("POST", endpoint, buffer)
	if err != nil {
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	tracer.InjectHTTPRequest(span, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("http status code[%d]", resp.StatusCode)
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	return nil
}
