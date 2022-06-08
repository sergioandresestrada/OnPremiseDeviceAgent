package queue

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQS defines the struct used to implement queue interface using AWS SQS
// It contains an SQS client and the queue URL to be used
type SQS struct {
	sqsClient *sqs.Client
	queueURL  *string
}

// NewQueueSQS creates and returns the reference to a new SQS struct
func NewQueueSQS() *SQS {
	q := &SQS{}
	q.initialize()
	return q
}

// SQSSendMessageAPI defines the interface for the GetQueueUrl and SendMessage functions.
// We use this interface to test the functions using a mocked service.
type SQSSendMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

func getQueueURL(c context.Context, api SQSSendMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func sendMsg(c context.Context, api SQSSendMessageAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(c, input)
}

func (queue *SQS) initialize() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	queue.sqsClient = sqs.NewFromConfig(cfg)

	queueName, ok := os.LookupEnv("SQS_QUEUE_NAME")
	if !ok {
		panic("Environment variable SQS_QUEUE_NAME does not exist.")
	}

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	result, err := getQueueURL(context.TODO(), queue.sqsClient, gQInput)
	if err != nil {
		panic(fmt.Sprintf("Got an error getting the queue URL: %v\n", err))
	}

	queue.queueURL = result.QueueUrl
}

// SendMessage receives an string and puts it in the correponding SQS URL
// Returns a non-nil error if there's one during the execution and nil otherwise
func (queue *SQS) SendMessage(s string) error {
	sMInput := &sqs.SendMessageInput{

		MessageBody:    aws.String(s),
		QueueUrl:       queue.queueURL,
		MessageGroupId: aws.String("1"),
	}

	resp, err := sendMsg(context.TODO(), queue.sqsClient, sMInput)
	if err != nil {
		err = fmt.Errorf("got an error sending the message to the queue: %w", err)
		return err
	}

	fmt.Println("Sent message with ID: " + *resp.MessageId)
	return nil
}
