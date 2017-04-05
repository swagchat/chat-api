package notification

import (
	"net/http"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

// [APNS] https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/PayloadKeyReference.html
// [FCM] https://firebase.google.com/docs/cloud-messaging/concept-options

type MessageInfo struct {
	Text  string
	Badge int64
}

type NotificationResult struct {
	Data          interface{}
	ProblemDetail *models.ProblemDetail
}

type NotificationChannel chan NotificationResult

type Provider interface {
	CreateTopic(string) NotificationChannel
	DeleteTopic(string) NotificationChannel
	CreateEndpoint(string) NotificationChannel
	DeleteEndpoint(string) NotificationChannel
	Subscribe(string, string) NotificationChannel
	Unsubscribe(string) NotificationChannel
	Publish(string, *MessageInfo) NotificationChannel
}

func GetProvider() Provider {
	var provider Provider
	switch utils.Cfg.ApiServer.Notification {
	case "awsSns":
		provider = &AwsSnsProvider{
			region:              utils.Cfg.AwsSns.Region,
			accessKeyId:         utils.Cfg.AwsSns.AccessKeyId,
			secretAccessKey:     utils.Cfg.AwsSns.SecretAccessKey,
			roomTopicNamePrefix: utils.Cfg.AwsSns.RoomTopicNamePrefix,
			applicationArn:      utils.Cfg.AwsSns.ApplicationArn,
		}
	}
	return provider
}

func createProblemDetail(title string, err error) *models.ProblemDetail {
	return &models.ProblemDetail{
		Title:     title,
		Status:    http.StatusInternalServerError,
		ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
		Detail:    err.Error(),
	}
}
