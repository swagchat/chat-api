package notification

import (
	"context"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type noopProvider struct {
	ctx context.Context
}

func (np *noopProvider) CreateTopic(roomId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	notificationTopicID := ""
	result := NotificationResult{
		Data: &notificationTopicID,
	}
	notificationChannel <- result
	return notificationChannel
}

func (np *noopProvider) DeleteTopic(notificationTopicId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *noopProvider) CreateEndpoint(userID string, platform scpb.Platform, token string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	notificationDeviceID := ""
	result := NotificationResult{
		Data: &notificationDeviceID,
	}
	notificationChannel <- result
	return notificationChannel
}

func (np *noopProvider) DeleteEndpoint(notificationDeviceId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *noopProvider) Subscribe(notificationTopicId string, notificationDeviceId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *noopProvider) Unsubscribe(notificationSubscribeId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *noopProvider) Publish(notificationTopicId, roomId string, messageInfo *MessageInfo) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}
