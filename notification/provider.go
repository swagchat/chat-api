package notification

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

// [APNS] https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/PayloadKeyReference.html
// [FCM] https://firebase.google.com/docs/cloud-messaging/concept-options

type MessageInfo struct {
	Text  string
	Badge int
}

type NotificationResult struct {
	Data          interface{}
	ProblemDetail *models.ProblemDetail
}

type NotificationChannel chan NotificationResult

type Provider interface {
	CreateTopic(string) NotificationChannel
	DeleteTopic(string) NotificationChannel
	CreateEndpoint(string, int, string) NotificationChannel
	DeleteEndpoint(string) NotificationChannel
	Subscribe(string, string) NotificationChannel
	Unsubscribe(string) NotificationChannel
	Publish(context.Context, string, string, *MessageInfo) NotificationChannel
}

func GetProvider() Provider {
	cfg := utils.GetConfig()

	var provider Provider
	switch cfg.Notification.Provider {
	case "awsSns":
		provider = &AwsSnsProvider{
			region:                cfg.Notification.AwsRegion,
			accessKeyId:           cfg.Notification.AwsAccessKeyId,
			secretAccessKey:       cfg.Notification.AwsSecretAccessKey,
			roomTopicNamePrefix:   cfg.Notification.RoomTopicNamePrefix,
			applicationArnIos:     cfg.Notification.AwsApplicationArnIos,
			applicationArnAndroid: cfg.Notification.AwsApplicationArnAndroid,
		}
	default:
		provider = &NotUseProvider{}
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
