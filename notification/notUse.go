package notification

type NotUseProvider struct {
}

func (provider NotUseProvider) CreateTopic(roomId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (provider NotUseProvider) DeleteTopic(notificationTopicId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (provider NotUseProvider) CreateEndpoint(deviceToken string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (provider NotUseProvider) DeleteEndpoint(notificationDeviceId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (provider NotUseProvider) Subscribe(notificationTopicId string, notificationDeviceId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (provider NotUseProvider) Unsubscribe(notificationSubscribeId string) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}

func (provider NotUseProvider) Publish(notificationTopicId string, messageInfo *MessageInfo) NotificationChannel {
	notificationChannel := make(NotificationChannel, 1)
	defer close(notificationChannel)
	result := NotificationResult{}
	notificationChannel <- result
	return notificationChannel
}
