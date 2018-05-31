package bots

import (
	"fmt"
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

type dialogFlowCredencial struct {
	ClientAccessToken string `json:"clientAccessToken"`
}

type dialogFlowResponse struct {
	ID        string           `json:"id"`
	Timestamp string           `json:"timestamp"`
	Lang      string           `json:"lang"`
	Result    dialogFlowResult `json:"result"`
}

type dialogFlowResult struct {
	Fulfillment dialogFlowFulfillment `json:"fulfillment"`
	Metadata    dialogFlowMetadata    `json:"metadata"`
	Score       float64               `json:"score"`
}

type dialogFlowMetadata struct {
	IntentID   string `json:"intentId"`
	IntentName string `json:"intentName"`
}

type dialogFlowFulfillment struct {
	Speech string
}

func (dp *dialogflowProvider) Post(m *models.Message, b *models.Bot, c utils.JSONText) *BotResult {
	var r BotResult

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

	var cred dialogFlowCredencial
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
	var res dialogFlowResponse

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Message: "Dialogflow response body unmarshal error",
			Error:   err,
		})
		return &r
	}

	score := res.Result.Score
	if res.Result.Metadata.IntentName == "Default Fallback Intent" {
		score = 0
	}

	var textPayload utils.JSONText
	payload := fmt.Sprintf("{\"text\":\"%s\",\"score\":%f}", res.Result.Fulfillment.Speech, score)
	err = json.Unmarshal([]byte(payload), &textPayload)
	general := models.RoleGeneral
	post := &models.Message{
		RoomID:    m.RoomID,
		UserID:    b.UserID,
		Type:      "textSuggest",
		Payload:   textPayload,
		EventName: "message",
		Role:      &general,
	}
	if b.Suggest {
		postCreated := m.Created + 1
		operator := models.RoleOperator

		post.SuggestMessageID = m.MessageID
		post.Role = &operator
		post.Created = postCreated
		post.Modified = postCreated
	}
	posts := make([]*models.Message, 0)
	posts = append(posts, post)
	messages := &models.Messages{
		Messages: posts,
	}
	r.Messages = messages

	return &r
}
