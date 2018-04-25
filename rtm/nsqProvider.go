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

type nsqProvider struct{}

func (np nsqProvider) Publish(rtmEvent *RTMEvent) error {
	cfg := utils.Config()
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(rtmEvent)

	url := utils.AppendStrings(cfg.RTM.NSQ.QueEndpoint, "/pub?topic=", cfg.RTM.NSQ.QueTopic)
	resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		return errors.Wrap(err, "NSQ post failure")
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "NSQ response body read failure")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(utils.AppendStrings("http status code[", strconv.Itoa(resp.StatusCode), "]"))
	}

	return nil
}
