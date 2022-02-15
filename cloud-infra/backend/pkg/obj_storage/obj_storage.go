package objstorage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

const (
	BUCKETNAME = "sergiotfgbucket"
)

var s3Client *s3.Client

func init() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	s3Client = s3.NewFromConfig(cfg)
}

func putFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func UploadFile(file *multipart.File, s3Name string) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(BUCKETNAME),
		Key:    aws.String(s3Name),
		Body:   *file,
	}

	_, err := putFile(context.TODO(), s3Client, input)
	if err != nil {
		fmt.Println("Got error uploading file:")
		fmt.Println(err)
	}

}
