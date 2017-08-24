package handlers

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var messengerHtmlData []byte = nil

func messengerHtmlHandler(rw http.ResponseWriter, req *http.Request) {
	if messengerHtmlData == nil {
		tmpExHtmlData, _ := ioutil.ReadFile("static/templates/messenger.html")
		tmpExHtml := string(tmpExHtmlData)
		tmpExHtml = strings.Replace(tmpExHtml, "SC_REACT_RTM_PROTOCOL", os.Getenv("SC_REACT_RTM_PROTOCOL"), 1)
		tmpExHtml = strings.Replace(tmpExHtml, "SC_REACT_RTM_HOST", os.Getenv("SC_REACT_RTM_HOST"), 1)
		tmpExHtml = strings.Replace(tmpExHtml, "SC_REACT_RTM_PATH", os.Getenv("SC_REACT_RTM_PATH"), 1)
		messengerHtmlData = []byte(tmpExHtml)
	}

	rw.Write(messengerHtmlData)
}
