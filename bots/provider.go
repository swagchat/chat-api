package bots

import (
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type MessageInfo struct {
	Text  string
	Badge int
}

type BotResult struct {
	Messages      *models.Messages
	ProblemDetail *models.ProblemDetail
}

type Provider interface {
	Post(*models.Message, *models.Bot, utils.JSONText) BotResult
}

func GetProvider(serviceName string) Provider {
	var provider Provider
	switch serviceName {
	case "A3RT":
		provider = &A3rtProvider{}
	case "API_AI":
		provider = &ApiAiProvider{}
	case "AWS_LEX":
		provider = &AwsLexProvider{}
	default:
		//provider = &NotUseProvider{}
	}
	return provider
}

func createProblemDetail(title string, err error) *models.ProblemDetail {
	//return &models.ProblemDetail{
	//	Title:     title,
	//	Status:    http.StatusInternalServerError,
	//	ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
	//	Detail:    err.Error(),
	//}
	return nil
}
