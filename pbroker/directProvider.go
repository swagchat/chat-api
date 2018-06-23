package pbroker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/utils"
)

type directProvider struct{}

func (dp directProvider) PublishMessage(rtmEvent *RTMEvent) error {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(rtmEvent)

	endpoint := fmt.Sprintf("%s/message", utils.Config().PBroker.Direct.Endpoint)
	resp, err := http.Post(
		endpoint,
		"application/json",
		buffer,
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("direct post failure. HTTP Endpoint=[%s]", endpoint))
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "direct response body read failure")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code[%d]", resp.StatusCode)
	}
	return nil
}
