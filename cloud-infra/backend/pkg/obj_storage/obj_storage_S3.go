package objstorage

import (
	"backend/pkg/types"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3 defines the struct used to implement ObjStorage interface using AWS S3
// It contains an S3 client and the name of the bucket to be used
type S3 struct {
	s3Client   *s3.Client
	BUCKETNAME string
	downloader *manager.Downloader
}

// NewObjStorageS3 creates and returns the reference to a new S3 struct
func NewObjStorageS3() *S3 {
	obj := &S3{}
	obj.initialize()
	return obj
}

// S3PutObjectAPI defines the interface for the PutObject function.
// We use this interface to test the function using a mocked service.
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// S3ListObjectsAPI defines the interface for the ListObjectsV2 function.
// We use this interface to test the function using a mocked service.
type S3ListObjectsAPI interface {
	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// GetObjects retrieves the objects in an Amazon Simple Storage Service (Amazon S3) bucket
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a ListObjectsV2Output object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to ListObjectsV2
func GetObjects(c context.Context, api S3ListObjectsAPI, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(c, input)
}

func (obj *S3) initialize() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic(fmt.Sprintf("Configuration error: %v\n", err))
	}

	_, ok := os.LookupEnv("S3_BUCKET_NAME")
	if !ok {
		panic("Environment variable S3_BUCKET_NAME does not exist")
	}
	obj.BUCKETNAME = os.Getenv("S3_BUCKET_NAME")

	obj.s3Client = s3.NewFromConfig(cfg)
	obj.downloader = manager.NewDownloader(obj.s3Client)
}

func putFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

// UploadFile receives an instance of a file that implements interface io.Reader and a name
// and upload that file with that name to S3
// Returns a non-nil error if there's one during the execution and nil otherwise
func (obj *S3) UploadFile(file io.Reader, s3Name string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(obj.BUCKETNAME),
		Key:    aws.String(s3Name),
		Body:   file,
	}

	_, err := putFile(context.TODO(), obj.s3Client, input)
	if err != nil {
		err = fmt.Errorf("got an error uploading the file: %w", err)
	}

	return err
}

// AvailableInformation returns an Information type object containing all the files (Jobs and Identification)
// whose name starts with 'Jobs-' or 'Identification-' respectively.
// It also returns a non-nil error if there's one during the execution and nil otherwise
func (obj *S3) AvailableInformation() (types.Information, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(obj.BUCKETNAME),
	}

	var listAvailable types.Information
	listAvailable.Jobs = make([]string, 0)
	listAvailable.Identification = make([]string, 0)

	resp, err := GetObjects(context.TODO(), obj.s3Client, input)

	if err != nil {
		err = fmt.Errorf("error getting the list of available files: %w", err)
		return listAvailable, err
	}

	for _, s := range resp.Contents {
		if strings.HasPrefix(*s.Key, "Jobs-") {
			listAvailable.Jobs = append(listAvailable.Jobs, *s.Key)
		}
		if strings.HasPrefix(*s.Key, "Identification-") {
			listAvailable.Identification = append(listAvailable.Identification, *s.Key)
		}
	}

	return listAvailable, nil
}

// GetFile receives a file name and a file pointer
// It will retrieve the mentioned file from S3 and store it in the pointer received
// It also returns a non-nil error if there's one during the execution and nil otherwise
func (obj *S3) GetFile(fileName string, fd *os.File) error {
	_, err := obj.downloader.Download(context.TODO(), fd, &s3.GetObjectInput{
		Bucket: aws.String(obj.BUCKETNAME),
		Key:    aws.String(fileName),
	})

	if err != nil {
		err = fmt.Errorf("error while downloading the file: %w", err)
	}

	return err
}
