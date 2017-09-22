package storage

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type AwsS3StorageProvider struct {
	accessKeyId        string
	secretAccessKey    string
	region             string
	acl                string
	uploadBucket       string
	uploadDirectory    string
	thumbnailBucket    string
	thumbnailDirectory string
}

func (provider AwsS3StorageProvider) Init() error {
	awsS3Client, err := provider.getSession()
	if err != nil {
		return err
	}

	params := &s3.CreateBucketInput{
		Bucket: aws.String(provider.uploadBucket),
	}
	createBucketResp, err := awsS3Client.CreateBucket(params)
	if err != nil {
		return err
	}
	log.Printf("Created bucket %#v", awsutil.StringValue(createBucketResp))

	params = &s3.CreateBucketInput{
		Bucket: aws.String(provider.thumbnailBucket),
	}
	createBucketResp, err = awsS3Client.CreateBucket(params)
	if err != nil {
		return err
	}
	log.Printf("Created bucket %#v", awsutil.StringValue(createBucketResp))
	return nil
}

func (provider AwsS3StorageProvider) Post(assetInfo *AssetInfo) (string, *models.ProblemDetail) {
	awsS3Client, err := provider.getSession()
	if err != nil {
		return "", &models.ProblemDetail{
			Title:     "Create session failed. (Amazon S3)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}

	byteData, err := ioutil.ReadAll(assetInfo.Data)
	if err != nil {
		return "", &models.ProblemDetail{
			Title:     "Reading asset data failed. (Amazon S3)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}

	data := bytes.NewReader(byteData)
	filePath := utils.AppendStrings(provider.uploadDirectory, "/", assetInfo.FileName)
	putObjectResp, err := awsS3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(provider.uploadBucket),
		Key:    aws.String(filePath),
		Body:   data,
		ACL:    &provider.acl,
	})
	if err != nil {
		return "", &models.ProblemDetail{
			Title:     "Create object failed. (Amazon S3)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}
	log.Println("Created object %#v", awsutil.StringValue(putObjectResp))

	sourceUrl := utils.AppendStrings("https://s3-ap-northeast-1.amazonaws.com/", provider.uploadBucket, "/", filePath)
	log.Println("sourceUrl:", sourceUrl)
	return sourceUrl, nil
}

func (provider AwsS3StorageProvider) Get(assetInfo *AssetInfo) ([]byte, *models.ProblemDetail) {
	return nil, nil
}

func (provider AwsS3StorageProvider) getSession() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(provider.region),
		Credentials: credentials.NewStaticCredentials(provider.accessKeyId, provider.secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}
