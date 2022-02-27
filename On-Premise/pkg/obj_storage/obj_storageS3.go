package obj_storage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	bucketName = "sergiotfgbucket"
)

type objStorageS3 struct {
	s3Client   *s3.Client
	downloader *manager.Downloader
}

func NewObjStorageS3() *objStorageS3 {
	obj := &objStorageS3{}
	obj.Initialize()
	return obj
}

func (obj *objStorageS3) Initialize() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("Object storage configuration error, " + err.Error())
	}

	obj.s3Client = s3.NewFromConfig(cfg)
	obj.downloader = manager.NewDownloader(obj.s3Client)
}

func (obj *objStorageS3) DownloadFile(message Message, fd *os.File) error {
	fmt.Printf("Downloading file %s\n", message.FileName)

	_, err := obj.downloader.Download(context.TODO(), fd, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(message.S3Name),
	})

	if err != nil {
		err = errors.New("Error while downloading the file: " + err.Error())
	}

	return err
}
