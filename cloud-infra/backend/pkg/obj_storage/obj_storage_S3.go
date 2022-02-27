package objstorage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	BUCKETNAME = "sergiotfgbucket"
)

type objStorageS3 struct {
	s3Client *s3.Client
}

func NewObjStorageS3() *objStorageS3 {
	obj := &objStorageS3{}
	obj.initialize()
	return obj
}

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func (obj *objStorageS3) initialize() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic(fmt.Sprintf("Configuration error: %v\n", err))
	}

	obj.s3Client = s3.NewFromConfig(cfg)
}

func putFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func (obj *objStorageS3) UploadFile(file *multipart.File, s3Name string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(BUCKETNAME),
		Key:    aws.String(s3Name),
		Body:   *file,
	}

	_, err := putFile(context.TODO(), obj.s3Client, input)
	if err != nil {
		err = fmt.Errorf("got an error uploading the file: %w", err)
	}

	return err
}
