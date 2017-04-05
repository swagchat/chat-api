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
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

type AwsS3StorageProvider struct {
	accessKeyId         string
	secretAccessKey     string
	region              string
	acl                 string
	userUploadBucket    string
	userUploadDirectory string
	thumbnailBucket     string
	thumbnailDirectory  string
}

var awsS3Client *s3.S3

func (provider AwsS3StorageProvider) Init() error {
	if awsS3Client == nil {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(utils.Cfg.AwsS3.Region),
			Credentials: credentials.NewStaticCredentials(utils.Cfg.AwsS3.AccessKeyId, utils.Cfg.AwsS3.SecretAccessKey, ""),
		})
		if err != nil {
			return err
		}
		awsS3Client := s3.New(sess)

		params := &s3.CreateBucketInput{
			Bucket: aws.String(utils.Cfg.AwsS3.UserUploadBucket),
		}
		createBucketResp, err := awsS3Client.CreateBucket(params)
		if err != nil {
			return err
		}
		log.Printf("Created bucket %#v", awsutil.StringValue(createBucketResp))

		params = &s3.CreateBucketInput{
			Bucket: aws.String(utils.Cfg.AwsS3.ThumbnailBucket),
		}
		createBucketResp, err = awsS3Client.CreateBucket(params)
		if err != nil {
			return err
		}
		log.Printf("Created bucket %#v", awsutil.StringValue(createBucketResp))
	}
	return nil
}

func (provider AwsS3StorageProvider) Post(assetInfo *AssetInfo) (string, *models.ProblemDetail) {
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
	filePath := utils.AppendStrings(provider.userUploadDirectory, "/", assetInfo.FileName)
	putObjectResp, err := awsS3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(provider.userUploadBucket),
		Key:    aws.String(filePath),
		Body:   data,
		ACL:    &utils.Cfg.AwsS3.Acl,
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

	sourceUrl := utils.AppendStrings("https://s3-ap-northeast-1.amazonaws.com/", provider.userUploadBucket, "/", filePath)
	log.Println("sourceUrl:", sourceUrl)
	return sourceUrl, nil
}

func (provider AwsS3StorageProvider) Get(assetInfo *AssetInfo) ([]byte, *models.ProblemDetail) {
	return nil, nil
}
