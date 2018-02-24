package utils

import (
	"errors"
	"fmt"
	"net/url"
)

func IsUrl(checkUrl string) error {
	urlStruct, err := url.Parse(checkUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("url parse error. %s", err.Error()))
	}
	schemes := []string{"http", "https"}
	for _, scheme := range schemes {
		if scheme == urlStruct.Scheme {
			return nil
		}
	}
	return errors.New("url is not http or https.")
}
