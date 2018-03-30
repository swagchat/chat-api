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

type NsqProvider struct{}

func (provider NsqProvider) Init() error {
	return nil
}

func (provider NsqProvider) PublishMessage(mi *MessagingInfo) error {
	rawIn := json.RawMessage(mi.Message)
	input, err := rawIn.MarshalJSON()
	cfg := utils.Config()

	url := utils.AppendStrings(cfg.RTM.NSQ.QueEndpoint, "/pub?topic=", cfg.RTM.NSQ.QueTopic)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(input))
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
