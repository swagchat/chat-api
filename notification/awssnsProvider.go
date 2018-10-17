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
	"github.com/betchi/tracer"
	logger "github.com/betchi/zapper"
	"github.com/swagchat/chat-api/config"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type awssnsProvider struct {
	ctx                   context.Context
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
	span := tracer.StartSpan(ap.ctx, "CreateTopic", "notification")
	defer tracer.Finish(span)

	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		params := &sns.CreateTopicInput{
			Name: aws.String(fmt.Sprintf("%s%s", config.Config().Notification.RoomTopicNamePrefix, roomId)),
		}
		createTopicOutput, err := client.CreateTopic(params)
		if err != nil {
			logger.Error(err.Error())
			tracer.SetError(span, err)
			result.Error = err
		} else {
			result.Data = createTopicOutput.TopicArn
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) DeleteTopic(notificationTopicId string) NotificationChannel {
	span := tracer.StartSpan(ap.ctx, "DeleteTopic", "notification")
	defer tracer.Finish(span)

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
			tracer.SetError(span, err)
			result.Error = err
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) CreateEndpoint(userID string, platform scpb.Platform, token string) NotificationChannel {
	span := tracer.StartSpan(ap.ctx, "CreateEndpoint", "notification")
	defer tracer.Finish(span)

	cfg := config.Config()

	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		var platformApplicationArn string
		switch platform {
		case scpb.Platform_PlatformIos:
			platformApplicationArn = cfg.Notification.AmazonSNS.ApplicationArnIos
		case scpb.Platform_PlatformAndroid:
			platformApplicationArn = cfg.Notification.AmazonSNS.ApplicationArnAndroid
		default:
			// TODO new error
			platformApplicationArn = ""
		}

		client := ap.newSnsClient()
		createPlatformEndpointInputParams := &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: aws.String(platformApplicationArn),
			Token:                  aws.String(token),
			CustomUserData:         aws.String(userID),
		}
		createPlatformEndpointOutput, err := client.CreatePlatformEndpoint(createPlatformEndpointInputParams)
		if err != nil {
			logger.Error(err.Error())
			tracer.SetError(span, err)
			result.Error = err
		} else {
			result.Data = createPlatformEndpointOutput.EndpointArn
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) DeleteEndpoint(notificationDeviceID string) NotificationChannel {
	span := tracer.StartSpan(ap.ctx, "DeleteEndpoint", "notification")
	defer tracer.Finish(span)

	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		deleteEndpointInputParams := &sns.DeleteEndpointInput{
			EndpointArn: aws.String(notificationDeviceID),
		}
		_, err := client.DeleteEndpoint(deleteEndpointInputParams)
		if err != nil {
			logger.Error(err.Error())
			tracer.SetError(span, err)
			result.Error = err
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) Subscribe(notificationTopicID string, notificationDeviceID string) NotificationChannel {
	span := tracer.StartSpan(ap.ctx, "Subscribe", "notification")
	defer tracer.Finish(span)

	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		subscribeInputParams := &sns.SubscribeInput{
			Protocol: aws.String("Application"),
			TopicArn: aws.String(notificationTopicID),
			Endpoint: aws.String(notificationDeviceID),
		}
		subscribeOutput, err := client.Subscribe(subscribeInputParams)
		if err != nil {
			logger.Error(err.Error())
			tracer.SetError(span, err)
			result.Error = err
		} else {
			result.Data = subscribeOutput.SubscriptionArn
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) Unsubscribe(notificationSubscribeID string) NotificationChannel {
	span := tracer.StartSpan(ap.ctx, "Unsubscribe", "notification")
	defer tracer.Finish(span)

	nc := make(NotificationChannel, 1)
	go func() {
		defer close(nc)
		result := NotificationResult{}

		client := ap.newSnsClient()
		params := &sns.UnsubscribeInput{
			SubscriptionArn: aws.String(notificationSubscribeID),
		}
		_, err := client.Unsubscribe(params)
		if err != nil {
			logger.Error(err.Error())
			tracer.SetError(span, err)
			result.Error = err
		}

		nc <- result
	}()
	return nc
}

func (ap *awssnsProvider) Publish(notificationTopicID, roomID string, messageInfo *MessageInfo) NotificationChannel {
	span := tracer.StartSpan(ap.ctx, "Publish", "notification")
	defer tracer.Finish(span)

	nc := make(NotificationChannel, 1)
	defer close(nc)
	result := NotificationResult{}

	client := ap.newSnsClient()
	contentAvailable := 1
	iosPush := iosPush{
		Alert:            messageInfo.Text,
		ContentAvailable: &contentAvailable,
		Sound:            "default",
		RoomId:           roomID,
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
		tracer.SetError(span, err)
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
		tracer.SetError(span, err)
		result.Error = err
		nc <- result
	}
	wrapper.GCM = string(b[:])
	pushData, err := json.Marshal(wrapper)
	if err != nil {
		logger.Error(err.Error())
		tracer.SetError(span, err)
		result.Error = err
		nc <- result
	}
	message := string(pushData[:])

	params := &sns.PublishInput{
		Message:          aws.String(message),
		MessageStructure: aws.String("json"),
		Subject:          aws.String("subject"),
		TopicArn:         aws.String(notificationTopicID),
	}
	res, err := client.Publish(params)
	if err != nil {
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil
	}
	logger.Info(fmt.Sprintf("[Amazon SNS]Publish message topicArn:%s message:%s response:%s", notificationTopicID, message, res.String()))

	nc <- result

	select {
	case <-ap.ctx.Done():
		return nc
	case <-nc:
		return nc
	}
}
