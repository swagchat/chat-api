package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/swagchat/chat-api/utils"
)

var messengerHTMLData []byte
var baseMessengerHTMLData = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>swagchat messenger</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/7.0.0/normalize.min.css" type="text/css" media="all">
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons" type="text/css" media="all">
    <link rel="stylesheet" href="http://localhost:9000/static/react-swagchat.min.css">
  </head>
  <body>
    <div id="swagchat" />
    <script src="http://localhost:9000/static/react-swagchat.min.js"></script>
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

func messengerHTMLHandler(rw http.ResponseWriter, req *http.Request) {
	if messengerHTMLData == nil {
		tmpExHTML := strings.Replace(baseMessengerHTMLData, "SC_REACT_RTM_PROTOCOL", os.Getenv("SC_REACT_RTM_PROTOCOL"), 1)
		tmpExHTML = strings.Replace(tmpExHTML, "SC_REACT_RTM_HOST", os.Getenv("SC_REACT_RTM_HOST"), 1)
		tmpExHTML = strings.Replace(tmpExHTML, "SC_REACT_RTM_PATH", os.Getenv("SC_REACT_RTM_PATH"), 1)

		chatEndpoint := os.Getenv("SC_REACT_CHAT_ENDPOINT")
		if chatEndpoint == "" {
			chatEndpoint = utils.AppendStrings("/", utils.API_VERSION)
		}
		tmpExHTML = strings.Replace(tmpExHTML, "SC_REACT_CHAT_ENDPOINT", chatEndpoint, 1)

		messengerHTMLData = []byte(tmpExHTML)
	}

	rw.Write(messengerHTMLData)
}
