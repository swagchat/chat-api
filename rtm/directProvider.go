package rtm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/utils"
)

type directProvider struct{}

func (dp directProvider) PublishMessage(mi *MessagingInfo) error {
	rawIn := json.RawMessage(mi.Message)
	input, err := rawIn.MarshalJSON()
	resp, err := http.Post(utils.AppendStrings(utils.Config().RTM.Direct.Endpoint, "/message"), "application/json", bytes.NewBuffer(input))
	if err != nil {
		return errors.Wrap(err, "Kafka post failure")
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Kafka response body read failure")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(utils.AppendStrings("http status code[", strconv.Itoa(resp.StatusCode), "]"))
	}
	return nil
}
