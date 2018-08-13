package pbroker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/tracer"
)

type directProvider struct {
	ctx context.Context
}

func (dp directProvider) PublishMessage(rtmEvent *RTMEvent) error {
	span := tracer.Provider(dp.ctx).StartSpan("PublishMessage", "pbroker")
	defer tracer.Provider(dp.ctx).Finish(span)

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(rtmEvent)

	endpoint := fmt.Sprintf("%s/message", config.Config().PBroker.Direct.Endpoint)

	req, err := http.NewRequest("POST", endpoint, buffer)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(dp.ctx).SetError(span, err)
		return err
	}

	tracer.Provider(dp.ctx).InjectHTTPRequest(span, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(dp.ctx).SetError(span, err)
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(dp.ctx).SetError(span, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("http status code[%d]", resp.StatusCode)
		logger.Error(err.Error())
		tracer.Provider(dp.ctx).SetError(span, err)
		return err
	}

	return nil
}
