package objstorage

import (
	"context"
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

// S3 defines the struct used to implement ObjStorage interface using AWS S3
// It contains an S3 client and the file downloader to be used
type S3 struct {
	s3Client   *s3.Client
	downloader *manager.Downloader
}

// NewObjStorageS3 creates and returns the reference to a new S3 struct
func NewObjStorageS3() *S3 {
	obj := &S3{}
	obj.initialize()
	return obj
}

func (obj *S3) initialize() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic(fmt.Sprintf("Object storage configuration error: %v", err))
	}

	obj.s3Client = s3.NewFromConfig(cfg)
	obj.downloader = manager.NewDownloader(obj.s3Client)
}

// DownloadFile downloads the file with name specified in received message and
// saves it to the given file pointer
// Returns a non-nil error if there's one during the execution and nil otherwise
func (obj *S3) DownloadFile(message Message, fd *os.File) error {
	fmt.Printf("Downloading file %s\n", message.FileName)

	_, err := obj.downloader.Download(context.TODO(), fd, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(message.S3Name),
	})

	if err != nil {
		err = fmt.Errorf("error while downloading the file: %w", err)
	}

	return err
}
