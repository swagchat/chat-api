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

type provider interface {
	Post(*models.Message, *models.Bot, utils.JSONText) BotResult
}

func Provider(serviceName string) provider {
	var p provider

	switch serviceName {
	case "A3RT":
		p = &a3rtProvider{}
	case "DIALOG_FLOW":
		p = &dialogflowProvider{}
	case "AWS_LEX":
		p = &awslexProvider{}
	}

	return p
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
