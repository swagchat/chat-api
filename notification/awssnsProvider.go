package notification

// [AWS SDK for Go Document] http://docs.aws.amazon.com/sdk-for-go/api/service/sns/

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

type awssnsProvider struct {
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

func (ap *awssnsProvider) newSnsClient() *sns.SNS {
	session, _ := session.NewSession(&aws.Config{
		Region:      aws.String(ap.region),
		Credentials: credentials.NewStaticCredentials(ap.accessKeyId, ap.secretAccessKey, ""),
	})
	return sns.New(session)
}

func (ap *awssnsProvider) CreateTopic(roomId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		params := &sns.CreateTopicInput{
			Name: aws.String(fmt.Sprintf("%s%s", utils.Config().Notification.RoomTopicNamePrefix, roomId)),
		}
		createTopicOutput, err := client.CreateTopic(params)
		if err != nil {
			logger.Error(err.Error())
			result.Error = err
		} else {
			result.Data = createTopicOutput.TopicArn
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) DeleteTopic(notificationTopicId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		params := &sns.DeleteTopicInput{
			TopicArn: aws.String(notificationTopicId),
		}
		_, err := client.DeleteTopic(params)
		if err != nil {
			logger.Error(err.Error())
			result.Error = err
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) CreateEndpoint(userId string, platform int32, deviceToken string) NotificationChannel {
	cfg := utils.Config()

	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		var platformApplicationArn string
		switch platform {
		case model.PlatformIOS:
			platformApplicationArn = cfg.Notification.AmazonSNS.ApplicationArnIos
		case model.PlatformAndroid:
			platformApplicationArn = cfg.Notification.AmazonSNS.ApplicationArnAndroid
		default:
			// TODO new error
			platformApplicationArn = ""
		}

		client := ap.newSnsClient()
		createPlatformEndpointInputParams := &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: aws.String(platformApplicationArn),
			Token:          aws.String(deviceToken),
			CustomUserData: aws.String(userId),
		}
		createPlatformEndpointOutput, err := client.CreatePlatformEndpoint(createPlatformEndpointInputParams)
		if err != nil {
			logger.Error(err.Error())
			result.Error = err
		} else {
			result.Data = createPlatformEndpointOutput.EndpointArn
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) DeleteEndpoint(notificationDeviceId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		deleteEndpointInputParams := &sns.DeleteEndpointInput{
			EndpointArn: aws.String(notificationDeviceId),
		}
		_, err := client.DeleteEndpoint(deleteEndpointInputParams)
		if err != nil {
			logger.Error(err.Error())
			result.Error = err
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) Subscribe(notificationTopicId string, notificationDeviceId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		subscribeInputParams := &sns.SubscribeInput{
			Protocol: aws.String("Application"),
			TopicArn: aws.String(notificationTopicId),
			Endpoint: aws.String(notificationDeviceId),
		}
		subscribeOutput, err := client.Subscribe(subscribeInputParams)
		if err != nil {
			logger.Error(err.Error())
			result.Error = err
		} else {
			result.Data = subscribeOutput.SubscriptionArn
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) Unsubscribe(notificationSubscribeId string) NotificationChannel {
	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		params := &sns.UnsubscribeInput{
			SubscriptionArn: aws.String(notificationSubscribeId),
		}
		_, err := client.Unsubscribe(params)
		if err != nil {
			logger.Error(err.Error())
			result.Error = err
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) Publish(ctx context.Context, notificationTopicId, roomId string, messageInfo *MessageInfo) NotificationChannel {
	nc := make(NotificationChannel, 1)
	defer close(nc)
	result := NotificationResult{}

	client := ap.newSnsClient()
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
		logger.Error(err.Error())
		result.Error = err
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
		logger.Error(err.Error())
		result.Error = err
		nc <- result
	}
	wrapper.GCM = string(b[:])
	pushData, err := json.Marshal(wrapper)
	if err != nil {
		logger.Error(err.Error())
		result.Error = err
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
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	logger.Info(fmt.Sprintf("[Amazon SNS]Publish message topicArn:%s message:%s response:%s", notificationTopicId, message, res.String()))

	nc <- result

	select {
	case <-ctx.Done():
		return nc
	case <-nc:
		return nc
	}
}
