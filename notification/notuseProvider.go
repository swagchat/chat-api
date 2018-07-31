package notification

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

type notuseProvider struct {
	ctx context.Context
}

func (np *notuseProvider) CreateTopic(roomId string) NotificationChannel {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "notification.awssnsProvider.CreateTopic")
	defer span.Finish()

	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) DeleteTopic(notificationTopicId string) NotificationChannel {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "notification.awssnsProvider.DeleteTopic")
	defer span.Finish()

	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) CreateEndpoint(userId string, platform int32, deviceToken string) NotificationChannel {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "notification.awssnsProvider.CreateEndpoint")
	defer span.Finish()

	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) DeleteEndpoint(notificationDeviceId string) NotificationChannel {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "notification.awssnsProvider.DeleteEndpoint")
	defer span.Finish()

	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) Subscribe(notificationTopicId string, notificationDeviceId string) NotificationChannel {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "notification.awssnsProvider.Subscribe")
	defer span.Finish()

	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) Unsubscribe(notificationSubscribeId string) NotificationChannel {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "notification.awssnsProvider.Unsubscribe")
	defer span.Finish()

	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (np *notuseProvider) Publish(notificationTopicId, roomId string, messageInfo *MessageInfo) NotificationChannel {
	span, _ := opentracing.StartSpanFromContext(np.ctx, "notification.awssnsProvider.Publish")
	defer span.Finish()

	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}
