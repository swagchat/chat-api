package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/tracer"
)

type awss3Provider struct {
	ctx                context.Context
	accessKeyId        string
	secretAccessKey    string
	region             string
	acl                string
	uploadBucket       string
	uploadDirectory    string
	thumbnailBucket    string
	thumbnailDirectory string
}

func (ap *awss3Provider) Init() error {
	span := tracer.Provider(ap.ctx).StartSpan("Init", "storage")
	defer tracer.Provider(ap.ctx).Finish(span)

	awsS3Client, err := ap.getSession()
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(ap.ctx).SetError(span, err)
		return err
	}

	params := &s3.CreateBucketInput{
		Bucket: aws.String(ap.uploadBucket),
	}
	_, err = awsS3Client.CreateBucket(params)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(ap.ctx).SetError(span, err)
		return err
	}

	params = &s3.CreateBucketInput{
		Bucket: aws.String(ap.thumbnailBucket),
	}
	_, err = awsS3Client.CreateBucket(params)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(ap.ctx).SetError(span, err)
		return err
	}

	return nil
}

func (ap *awss3Provider) Post(assetInfo *AssetInfo) (string, error) {
	span := tracer.Provider(ap.ctx).StartSpan("Post", "storage")
	defer tracer.Provider(ap.ctx).Finish(span)

	awsS3Client, err := ap.getSession()
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(ap.ctx).SetError(span, err)
		return "", err
	}

	byteData, err := ioutil.ReadAll(assetInfo.Data)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(ap.ctx).SetError(span, err)
		return "", err
	}

	data := bytes.NewReader(byteData)
	filePath := fmt.Sprintf("%s/%s", ap.uploadDirectory, assetInfo.Filename)
	_, err = awsS3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(ap.uploadBucket),
		Key:    aws.String(filePath),
		Body:   data,
		ACL:    &ap.acl,
	})
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(ap.ctx).SetError(span, err)
		return "", err
	}

	sourceURL := fmt.Sprintf("https://s3-ap-northeast-1.amazonaws.com/%s/%s", ap.uploadBucket, filePath)
	return sourceURL, nil
}

func (ap *awss3Provider) Get(assetInfo *AssetInfo) ([]byte, error) {
	span := tracer.Provider(ap.ctx).StartSpan("Get", "storage")
	defer tracer.Provider(ap.ctx).Finish(span)

	return nil, nil
}

func (ap *awss3Provider) getSession() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(ap.region),
		Credentials: credentials.NewStaticCredentials(ap.accessKeyId, ap.secretAccessKey, ""),
	})
	if err != nil {
		return nil, errors.Wrap(err, "AWS S3 create session failure")
	}
	return s3.New(sess), nil
}
