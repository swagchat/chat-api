package handlers

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
    <script src="SWAG_API_ENDPOINT/static/SWAG_JS_FILE_NAME"></script>
    <script>
      Swag.renderMessenger({
        clientParams: {
          apiEndpoint: 'SWAG_API_ENDPOINT',
          wsEndpoint: 'SWAG_WS_ENDPOINT',
          userId: 'SWAG_USER_ID',
          username: 'SWAG_USER_NAME',
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
		tmpExHTML := strings.Replace(baseMessengerHTMLData, "SWAG_API_ENDPOINT", os.Getenv("SWAG_API_ENDPOINT"), -1)
		tmpExHTML = strings.Replace(tmpExHTML, "SWAG_WS_ENDPOINT", os.Getenv("SWAG_WS_ENDPOINT"), -1)
		tmpExHTML = strings.Replace(tmpExHTML, "SWAG_USER_ID", os.Getenv("SWAG_USER_ID"), -1)
		tmpExHTML = strings.Replace(tmpExHTML, "SWAG_USER_NAME", os.Getenv("SWAG_USER_NAME"), -1)
		tmpExHTML = strings.Replace(tmpExHTML, "SWAG_JS_FILE_NAME", os.Getenv("SWAG_JS_FILE_NAME"), -1)

		messengerHTMLData = []byte(tmpExHTML)
	}

	rw.Write(messengerHTMLData)
}
