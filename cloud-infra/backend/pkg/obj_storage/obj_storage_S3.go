package objstorage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3 defines the struct used to implement ObjStorage interface using AWS S3
// It contains an S3 client and the name of the bucket to be used
type S3 struct {
	s3Client   *s3.Client
	BUCKETNAME string
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
