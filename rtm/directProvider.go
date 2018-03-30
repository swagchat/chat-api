package rtm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/swagchat/chat-api/utils"
)

type DirectProvider struct{}

func (provider DirectProvider) Init() error {
	return nil
}

func (provider DirectProvider) PublishMessage(mi *MessagingInfo) error {
	rawIn := json.RawMessage(mi.Message)
	input, err := rawIn.MarshalJSON()
	resp, err := http.Post(utils.AppendStrings(utils.GetConfig().RTM.Direct.Endpoint, "/message"), "application/json", bytes.NewBuffer(input))
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(utils.AppendStrings("http status code[", strconv.Itoa(resp.StatusCode), "]"))
	}
	return nil
}
