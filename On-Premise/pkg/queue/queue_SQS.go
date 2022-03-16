package queue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	queueName = "test.fifo"
	waitTime  = 20
)

// SQS defines the struct used to implement queue interface using AWS SQS
// It contains an SQS client, the queue URL to be used, and the input struct to receive messages
type SQS struct {
	sqsClient *sqs.Client
	queueURL  *string
	mInput    *sqs.ReceiveMessageInput
}

// NewQueueSQS creates and returns the reference to a new SQS struct
func NewQueueSQS() *SQS {
	q := &SQS{}
	q.initialize()
	return q
}

// SQSGetLPMsgAPI defines the interface for the GetQueueUrl and ReceiveMessage functions.
// We use this interface to test the functions using a mocked service.
type SQSGetLPMsgAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

// SQSDeleteMessageAPI defines the interface for the GetQueueUrl and DeleteMessage functions.
// We use this interface to test the functions using a mocked service.
type SQSDeleteMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

func getQueueURL(c context.Context, api SQSGetLPMsgAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func getLPMessages(c context.Context, api SQSGetLPMsgAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

func removeMessage(c context.Context, api SQSDeleteMessageAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}

func (queue *SQS) initialize() {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic(fmt.Sprintf("configuration error: %v", err))
	}

	queue.sqsClient = sqs.NewFromConfig(cfg)

	queueNameString := aws.String(queueName)

	qInput := &sqs.GetQueueUrlInput{
		QueueName: queueNameString,
	}

	result, err := getQueueURL(context.TODO(), queue.sqsClient, qInput)
	if err != nil {
		panic(fmt.Sprintf("Got an error getting the queue URL: %v", err))
	}

	queue.queueURL = result.QueueUrl

	queue.mInput = &sqs.ReceiveMessageInput{
		QueueUrl: queue.queueURL,
		AttributeNames: []types.QueueAttributeName{
			"SentTimestamp",
		},
		MaxNumberOfMessages: 1,
		MessageAttributeNames: []string{
			"All",
		},
		WaitTimeSeconds: int32(waitTime),
	}

}

// ReceiveMessages used the queue to retrieve and return Messages from it
// Returns nil if there's an error receiving messages
func (queue *SQS) ReceiveMessages() []types.Message {

	resp, err := getLPMessages(context.TODO(), queue.sqsClient, queue.mInput)

	if err != nil {
		fmt.Printf("Got an error receiving messages: %v\n", err)
		return nil
	}

	return resp.Messages
}

// RemoveMessage received a processed message and removes it from the queue
// Returns a non-nil error if there's one during the execution and nil otherwise
func (queue *SQS) RemoveMessage(msg types.Message) error {

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      queue.queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	}

	_, err := removeMessage(context.TODO(), queue.sqsClient, dMInput)

	if err != nil {
		err = fmt.Errorf("got an error deleting the message fron the queue: %w", err)
	}

	return err

}
