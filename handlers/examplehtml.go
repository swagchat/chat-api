package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/fairway-corp/swagchat-api/utils"
)

var messengerHtmlData []byte
var baseMessengerHtmlData = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>swagchat messenger</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/7.0.0/normalize.min.css" type="text/css" media="all">
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons" type="text/css" media="all">
    <link rel="stylesheet" href="https://unpkg.com/react-swagchat/dist/react-swagchat.min.css">
  </head>
  <body>
    <div id="swagchat" />
    <script src="https://unpkg.com/react-swagchat/dist/react-swagchat.min.js"></script>
    <script>
      Swag.renderMessenger({
        userId: '00581ea9-3547-4c81-930c-a3ed042e4b21',
        apiEndpoint: 'SC_REACT_CHAT_ENDPOINT',
        rtmProtocol: 'SC_REACT_RTM_PROTOCOL',
        rtmHost: 'SC_REACT_RTM_HOST',
        rtmPath: 'SC_REACT_RTM_PATH',
      });
    </script>
  </body>
</html>`

func messengerHtmlHandler(rw http.ResponseWriter, req *http.Request) {
	if messengerHtmlData == nil {
		tmpExHtml := strings.Replace(baseMessengerHtmlData, "SC_REACT_RTM_PROTOCOL", os.Getenv("SC_REACT_RTM_PROTOCOL"), 1)
		tmpExHtml = strings.Replace(tmpExHtml, "SC_REACT_RTM_HOST", os.Getenv("SC_REACT_RTM_HOST"), 1)
		tmpExHtml = strings.Replace(tmpExHtml, "SC_REACT_RTM_PATH", os.Getenv("SC_REACT_RTM_PATH"), 1)

		chatEndpoint := os.Getenv("SC_REACT_CHAT_ENDPOINT")
		if chatEndpoint == "" {
			chatEndpoint = utils.AppendStrings("/", utils.API_VERSION)
		}
		tmpExHtml = strings.Replace(tmpExHtml, "SC_REACT_CHAT_ENDPOINT", chatEndpoint, 1)

		messengerHtmlData = []byte(tmpExHtml)
	}

	rw.Write(messengerHtmlData)
}
