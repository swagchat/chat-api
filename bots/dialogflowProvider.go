package bots

import (
	"net/http"
	"net/url"

	"encoding/json"
	"io/ioutil"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

type dialogflowProvider struct {
}

func (dp *dialogflowProvider) Post(m *models.Message, b *models.Bot, c utils.JSONText) BotResult {
	r := BotResult{}

	var message string
	switch m.Type {
	case "text":
		var payloadText models.PayloadText
		json.Unmarshal(m.Payload, &payloadText)
		message = payloadText.Text
	case "image":
		message = "画像を受信しました"
	default:
		message = "メッセージを受信しました"
	}

	var cred models.ApiAiCredencial
	json.Unmarshal(c, &cred)

	values := url.Values{}
	values.Set("v", "20150910")
	values.Add("timezone", "Asia/Tokyo")
	values.Add("lang", "ja")
	values.Add("sessionId", b.UserID)
	values.Add("query", message)

	req, err := http.NewRequest(
		"GET",
		"https://api.api.ai/v1/query?"+values.Encode(),
		nil,
	)
	if err != nil {
	}
	req.Header.Set("Authorization", utils.AppendStrings("Bearer ", cred.ClientAccessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	var res models.ApiAiResponse

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Message: "Dialogflow response body unmarshal error",
			Error:   err,
		})
		return r
	}

	var textPayload json.RawMessage
	err = json.Unmarshal([]byte("{\"text\": \""+res.Result.Fulfillment.Speech+"\"}"), &textPayload)
	post := &models.Message{
		RoomID:    m.RoomID,
		UserID:    b.UserID,
		Type:      "text",
		Payload:   textPayload,
		EventName: "message",
	}
	posts := make([]*models.Message, 0)
	posts = append(posts, post)
	messages := &models.Messages{
		Messages: posts,
	}
	r.Messages = messages

	return r
}
