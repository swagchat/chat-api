package bots

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexruntimeservice"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type AwsLexProvider struct {
}

func (p *AwsLexProvider) Post(m *models.Message, b *models.Bot, c utils.JSONText) BotResult {
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
		UserId: aws.String(m.UserId),
	}
	output, err := svc.PostText(input)
	if err != nil {
		log.Println(err)
	}

	log.Printf("============= Amazon Lex ====================")
	log.Printf("%#v\n", output)

	var textPayload utils.JSONText
	err = json.Unmarshal([]byte("{\"text\": \""+*output.Message+"\"}"), &textPayload)
	post := &models.Message{
		RoomId:  m.RoomId,
		UserId:  b.UserId,
		Type:    "text",
		Payload: textPayload,
	}
	posts := make([]*models.Message, 0)
	posts = append(posts, post)
	messages := &models.Messages{
		Messages: posts,
	}
	r.Messages = messages

	return r
}
