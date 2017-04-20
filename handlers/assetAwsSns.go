package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/fairway-corp/swagchat-api/utils"
)

func SetAssetAwsSnsMux() {
	Mux.PostFunc("/assets/aws-sns", ColsHandler(PostAssetAwsSns))
}

type AwsSNSSubscribeInput struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
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

type AssetS3SNSRecords struct {
	Records []Record `json:"Records"`
}

type Record struct {
	S3 S3 `json:"s3"`
}

type S3 struct {
	Object Object `json:"object"`
}

type Object struct {
	Etag      string `json:"eTag"`
	Key       string `json:"key"`
	Sequencer string `json:"sequencer"`
	Size      int    `json:"size"`
}

func PostAssetAwsSns(w http.ResponseWriter, r *http.Request) {
	log.Println("PostAssetSNS")

	decoder := json.NewDecoder(r.Body)
	var input AwsSNSSubscribeInput
	err := decoder.Decode(&input)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%#v", input)
	if input.Type == "SubscriptionConfirmation" {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(utils.Cfg.AwsS3.Region),
			Credentials: credentials.NewStaticCredentials(utils.Cfg.AwsS3.AccessKeyId, utils.Cfg.AwsS3.SecretAccessKey, ""),
		})
		cli := sns.New(sess)

		params := &sns.ConfirmSubscriptionInput{
			Token:    aws.String(input.Token),
			TopicArn: aws.String(input.TopicArn),
		}
		resp, err := cli.ConfirmSubscription(params)

		if err != nil {
			log.Println(err.Error())
			return
		}

		log.Println(resp)
	} else {
		log.Println(input.Message)
		var records AssetS3SNSRecords
		err = json.Unmarshal([]byte(input.Message), &records)
		if err != nil {
			log.Println(err)
			return
		}
		filePath := records.Records[0].S3.Object.Key
		filename := filepath.Base(filePath)
		log.Println(filePath)
		log.Println(filename)
		pos := strings.LastIndex(filename, ".")
		messageId := filename[0:pos]
		log.Println("messageId=", messageId)

		// TODO Update PayloadImage.ThumbnailUrl
		// https://s3-ap-northeast-1.amazonaws.com/swagchat-thumbnail/assets/reduced/2d2eac8b-3ff2-491a-9ae4-e73009caecd6.png
		/*
			url := "https://s3-" + utils.Cfg.AwsS3.Region + ".amazonaws.com/" + utils.Cfg.AwsS3.ThumbnailBucketNameAsset + "/" + filePath
			log.Println("url=", url)

			if strings.Contains(filePath, utils.Cfg.AwsS3.SourceDirectory) {
				log.Println("SourceDirectory")
			} else if strings.Contains(filePath, utils.Cfg.AwsS3.ThumbnailDirectory) {
				log.Println("ThumbnailDirectory")

				message, problemDetail := services.GetMessage(messageId)
				if problemDetail != nil {
					log.Println(problemDetail)
					return
				}
				var payloadImage models.PayloadImage
				json.Unmarshal(message.Payload, &payloadImage)
				payloadImage.ThumbnailUrl = url
				payloadImageBuffer, err := json.Marshal(payloadImage)
				message.Payload = payloadImageBuffer
				datastoreProvider := datastore.GetDatastoreProvider()
				<-datastoreProvider.MessageUpdate(message)
			}
		*/
	}
}
