package bots

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexruntimeservice"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

type awslexProvider struct {
}

func (ap *awslexProvider) Post(m *models.Message, b *models.Bot, c utils.JSONText) BotResult {
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

	var cred models.AwsLexCredencial
	json.Unmarshal(c, &cred)

	session, _ := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(cred.AccessKeyID, cred.SecretAccessKey, ""),
	})
	svc := lexruntimeservice.New(session)

	input := &lexruntimeservice.PostTextInput{
		BotName:   aws.String("BookTrip"),
		BotAlias:  aws.String("lextest"),
		InputText: aws.String(message),
		//SessionAttributes: "",
		UserId: aws.String(m.UserID),
	}
	output, err := svc.PostText(input)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Message: "Amazon Lex post error",
			Error:   err,
		})
		return r
	}

	var textPayload utils.JSONText
	err = json.Unmarshal([]byte("{\"text\": \""+*output.Message+"\"}"), &textPayload)
	post := &models.Message{
		RoomID:  m.RoomID,
		UserID:  b.UserID,
		Type:    "text",
		Payload: textPayload,
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
