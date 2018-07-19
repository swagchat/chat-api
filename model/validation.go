package model

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

// IsURL is url
func IsURL(checkURL string) error {
	urlStruct, err := url.Parse(checkURL)
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

// IsValidID is valid ID
func IsValidID(ID string) bool {
	r := regexp.MustCompile(`(?m)^[0-9a-z-]+$`)
	return r.MatchString(ID)
}
