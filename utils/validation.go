package utils

import (
	"errors"
	"net/url"
)

func IsUrl(checkUrl string) error {
	urlStruct, err := url.Parse(checkUrl)
	if err != nil {
		return errors.New(AppendStrings("url parse error. ", err.Error()))
	}
	schemes := []string{"http", "https"}
	for _, scheme := range schemes {
		if scheme == urlStruct.Scheme {
			return nil
		}
	}
	return errors.New("url is not http or https.")
}
