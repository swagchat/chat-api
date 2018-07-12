package notification

import "context"

type notuseProvider struct {
}

func (np *notuseProvider) CreateTopic(roomId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) DeleteTopic(notificationTopicId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) CreateEndpoint(userId string, platform int32, deviceToken string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) DeleteEndpoint(notificationDeviceId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) Subscribe(notificationTopicId string, notificationDeviceId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) Unsubscribe(notificationSubscribeId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) Publish(ctx context.Context, notificationTopicId, roomId string, messageInfo *MessageInfo) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}
