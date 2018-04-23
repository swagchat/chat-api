package rtm

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	endpoint := fmt.Sprintf("%s/message", utils.Config().RTM.Direct.Endpoint)
	resp, err := http.Post(
		endpoint,
		"application/json",
		bytes.NewBuffer(input),
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("direct post failure. HTTP Endpoint=[%s]", endpoint))
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "direct response body read failure")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(utils.AppendStrings("http status code[", strconv.Itoa(resp.StatusCode), "]"))
	}
	return nil
}
