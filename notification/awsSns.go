package notification

// [AWS SDK for Go Document] http://docs.aws.amazon.com/sdk-for-go/api/service/sns/

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type AwsSnsProvider struct {
	region                string
	accessKeyId           string
	secretAccessKey       string
	roomTopicNamePrefix   string
	applicationArnIos     string
	applicationArnAndroid string
}

type PushNotification struct {
	sns *sns.SNS
}

type wrapper struct {
	APNS        string `json:"APNS"`
	APNSSandbox string `json:"APNS_SANDBOX"`
	Default     string `json:"default"`
	GCM         string `json:"GCM"`
}

type iosPushWrapper struct {
	APS iosPush `json:"aps"`
}

type iosPush struct {
	Alert            string `json:"alert,omitempty"`
	Badge            *int   `json:"badge,omitempty"`
	Sound            string `json:"sound,omitempty"`
	ContentAvailable *int   `json:"content-available,omitempty"`
	Category         string `json:"category,omitempty"`
	ThreadId         string `json:"thread-id,omitempty"`
	RoomId           string `json:"roomId,omitempty"`
}

type gcmPushWrapper struct {
	Data gcmPush `json:"data"`
}

type gcmPush struct {
	Message string      `json:"message,omitempty"`
	Custom  interface{} `json:"custom"`
	Badge   *int        `json:"badge,omitempty"`
}

func (provider AwsSnsProvider) Init() error {
	return nil
}

func (provider AwsSnsProvider) newSnsClient() *sns.SNS {
	session, _ := session.NewSession(&aws.Config{
		Region:      aws.String(provider.region),
		Credentials: credentials.NewStaticCredentials(provider.accessKeyId, provider.secretAccessKey, ""),
	})
	return sns.New(session)
}

func (provider AwsSnsProvider) CreateTopic(roomId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := provider.newSnsClient()
		params := &sns.CreateTopicInput{
			Name: aws.String(utils.AppendStrings(utils.Cfg.Notification.RoomTopicNamePrefix, roomId)),
		}
		createTopicOutput, err := client.CreateTopic(params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating topic.", err)
		} else {
			result.Data = createTopicOutput.TopicArn
		}

		nc <- result
	}()
	return nc
}

func (provider AwsSnsProvider) DeleteTopic(notificationTopicId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := provider.newSnsClient()
		params := &sns.DeleteTopicInput{
			TopicArn: aws.String(notificationTopicId),
		}
		_, err := client.DeleteTopic(params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while deleting topic.", err)
		}

		nc <- result
	}()
	return nc
}

func (provider AwsSnsProvider) CreateEndpoint(userId string, platform int, deviceToken string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		var platformApplicationArn string
		switch platform {
		case models.PLATFORM_IOS:
			platformApplicationArn = utils.Cfg.Notification.AwsApplicationArnIos
		case models.PLATFORM_ANDROID:
			platformApplicationArn = utils.Cfg.Notification.AwsApplicationArnAndroid
		default:
			// TODO new error
			platformApplicationArn = ""
		}

		client := provider.newSnsClient()
		createPlatformEndpointInputParams := &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: aws.String(platformApplicationArn),
			Token:          aws.String(deviceToken),
			CustomUserData: aws.String(userId),
		}
		createPlatformEndpointOutput, err := client.CreatePlatformEndpoint(createPlatformEndpointInputParams)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating endpoint.", err)
		} else {
			result.Data = createPlatformEndpointOutput.EndpointArn
		}

		nc <- result
	}()
	return nc
}

func (provider AwsSnsProvider) DeleteEndpoint(notificationDeviceId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := provider.newSnsClient()
		deleteEndpointInputParams := &sns.DeleteEndpointInput{
			EndpointArn: aws.String(notificationDeviceId),
		}
		_, err := client.DeleteEndpoint(deleteEndpointInputParams)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while deleting endpoint.", err)
		}

		nc <- result
	}()
	return nc
}

func (provider AwsSnsProvider) Subscribe(notificationTopicId string, notificationDeviceId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := provider.newSnsClient()
		subscribeInputParams := &sns.SubscribeInput{
			Protocol: aws.String("Application"),
			TopicArn: aws.String(notificationTopicId),
			Endpoint: aws.String(notificationDeviceId),
		}
		subscribeOutput, err := client.Subscribe(subscribeInputParams)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while subscribing.", err)
		} else {
			result.Data = subscribeOutput.SubscriptionArn
		}

		nc <- result
	}()
	return nc
}

func (provider AwsSnsProvider) Unsubscribe(notificationSubscribeId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := provider.newSnsClient()
		params := &sns.UnsubscribeInput{
			SubscriptionArn: aws.String(notificationSubscribeId),
		}
		_, err := client.Unsubscribe(params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while unsubscribing.", err)
		}

		nc <- result
	}()
	return nc
}

func (provider AwsSnsProvider) Publish(ctx context.Context, notificationTopicId, roomId string, messageInfo *MessageInfo) NotificationChannel {
	nc := make(NotificationChannel, 1)
	defer close(nc)
	result := NotificationResult{}

	client := provider.newSnsClient()
	contentAvailable := 1
	iosPush := iosPush{
		Alert:            messageInfo.Text,
		ContentAvailable: &contentAvailable,
		Sound:            "default",
		RoomId:           roomId,
	}
	if &messageInfo.Badge != nil {
		iosPush.Badge = &messageInfo.Badge
	}

	wrapper := wrapper{}
	ios := iosPushWrapper{
		APS: iosPush,
	}
	b, err := json.Marshal(ios)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while publishing.", err)
		nc <- result
	}
	wrapper.APNS = string(b[:])
	wrapper.APNSSandbox = wrapper.APNS
	wrapper.Default = messageInfo.Text
	gcm := gcmPushWrapper{
		Data: gcmPush{
			Message: messageInfo.Text,
			Badge:   &messageInfo.Badge,
		},
	}
	b, err = json.Marshal(gcm)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while publishing.", err)
		nc <- result
	}
	wrapper.GCM = string(b[:])
	pushData, err := json.Marshal(wrapper)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while publishing.", err)
		nc <- result
	}
	message := string(pushData[:])

	params := &sns.PublishInput{
		Message:          aws.String(message),
		MessageStructure: aws.String("json"),
		Subject:          aws.String("subject"),
		TopicArn:         aws.String(notificationTopicId),
	}
	res, err := client.Publish(params)
	utils.AppLogger.Info("",
		zap.String("msg", "[Amazon SNS]Publish message."),
		zap.String("topicArn", notificationTopicId),
		zap.String("message", message),
		zap.String("response", res.String()),
	)
	nc <- result

	select {
	case <-ctx.Done():
		return nc
	case <-nc:
		return nc
	}
}
