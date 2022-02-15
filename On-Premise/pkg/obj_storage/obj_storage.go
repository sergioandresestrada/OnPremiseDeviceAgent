package obj_storage

import (
	"context"
	"fmt"
	"os"

	"On-Premise/pkg/types"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type Message = types.Message

const (
	bucketName = "sergiotfgbucket"
)

var s3client *s3.Client
var downloader *manager.Downloader

func init() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("Object storage configuration error, " + err.Error())
	}

	s3client = s3.NewFromConfig(cfg)
	downloader = manager.NewDownloader(s3client)

}

func DownloadFile(message Message, fd *os.File) {
	fmt.Printf("Downloading file %s\n", message.FileName)

	_, err := downloader.Download(context.TODO(), fd, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(message.S3Name),
	})

	if err != nil {
		fmt.Println("Error while downloading the file")
	}
}
