package rest

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/tracer"
)

func setAssetAwsSnsMux() {
	mux.PostFunc("/assets/aws-sns", commonHandler(postAssetAwsSns))
}

// AwsSNSSubscribeInput is AwsSNSSubscribeInput
type AwsSNSSubscribeInput struct {
	Type             string `json:"Type"`
	MessageID        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Subject          string `json:"Subject"`
	Message          string `json:"Message,omitempty"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	SubscribeURL     string `json:"SubscribeURL,omitempty"`
	UnsubscribeURL   string `json:"UnsubscribeURL,omitempty"`
	Token            string `json:"Token,omitempty"`
}

// AssetS3SNSRecords is asset s3 sns records
type AssetS3SNSRecords struct {
	Records []Record `json:"Records"`
}

// Record is Record
type Record struct {
	S3 S3 `json:"s3"`
}

// S3 is s3
type S3 struct {
	Object Object `json:"object"`
}

// Object is Object
type Object struct {
	Etag      string `json:"eTag"`
	Key       string `json:"key"`
	Sequencer string `json:"sequencer"`
	Size      int    `json:"size"`
}

func postAssetAwsSns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("postAssetAwsSns", "rest")
	defer tracer.Provider(ctx).Finish(span)

	decoder := json.NewDecoder(r.Body)
	var input AwsSNSSubscribeInput
	err := decoder.Decode(&input)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if input.Type == "SubscriptionConfirmation" {
		cfg := config.Config()
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(cfg.Storage.AWSS3.Region),
			Credentials: credentials.NewStaticCredentials(cfg.Storage.AWSS3.AccessKeyID, cfg.Storage.AWSS3.SecretAccessKey, ""),
		})
		cli := sns.New(sess)

		params := &sns.ConfirmSubscriptionInput{
			Token:    aws.String(input.Token),
			TopicArn: aws.String(input.TopicArn),
		}
		_, err = cli.ConfirmSubscription(params)

		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		var records AssetS3SNSRecords
		err = json.Unmarshal([]byte(input.Message), &records)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		// filePath := records.Records[0].S3.Object.Key
		// filename := filepath.Base(filePath)
		// pos := strings.LastIndex(filename, ".")
		// messageId := filename[0:pos]
	}
}
