package handler

import (
	"net/http"
	"os"
	"strings"
)

var messengerHTMLData []byte
var baseMessengerHTMLData = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>swagchat messenger</title>
    <meta name="viewport" content="width=device-width,initial-scale=1,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/7.0.0/normalize.min.css" type="text/css" media="all">
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons" type="text/css" media="all">
  </head>
  <body>
    <div id="swag" />
    <script src="http://localhost:8101/static/main.5785e4e6.js"></script>
    <script>
      Swag.renderMessenger({
        clientParams: {
          apiEndpoint: 'http://localhost:8101',
          wsEndpoint: 'ws://localhost:8102',
          userId: '958c775a-9d71-4b26-8b17-5a210926a75e',
          username: 'demo user',
          paths: {
            roomListPath: '/rooms',
          }
        }
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
			chatEndpoint = "/"
		}
		tmpExHTML = strings.Replace(tmpExHTML, "SC_REACT_CHAT_ENDPOINT", chatEndpoint, 1)

		messengerHTMLData = []byte(tmpExHTML)
	}

	rw.Write(messengerHTMLData)
}
