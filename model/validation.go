package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/swagchat/chat-api/logger"
)

func isURL(checkURL string) error {
	urlStruct, err := url.Parse(checkURL)
	if err != nil {
		err := fmt.Errorf("url parse error. %s", err.Error())
		logger.Error(err.Error())
		return err
	}
	schemes := []string{"http", "https"}
	for _, scheme := range schemes {
		if scheme == urlStruct.Scheme {
			return nil
		}
	}

	err = fmt.Errorf("url is not http or https")
	logger.Error(err.Error())
	return err
}

func isValidID(ID string) bool {
	r := regexp.MustCompile(`(?m)^[0-9a-z-]+$`)
	return r.MatchString(ID)
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
