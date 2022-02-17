package queue

import (
	"context"
	"errors"
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

type queueSQS struct {
	sqsClient *sqs.Client
	queueURL  *string
	mInput    *sqs.ReceiveMessageInput
}

func NewQueueSQS() *queueSQS {
	q := &queueSQS{}
	q.initialize()
	return q
}

type SQSGetLPMsgAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

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

func (queue *queueSQS) initialize() {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	queue.sqsClient = sqs.NewFromConfig(cfg)

	queueNameString := aws.String(queueName)

	qInput := &sqs.GetQueueUrlInput{
		QueueName: queueNameString,
	}

	result, err := getQueueURL(context.TODO(), queue.sqsClient, qInput)
	if err != nil {
		panic("Got an error getting the queue URL: " + err.Error())
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

func (queue *queueSQS) ReceiveMessages() []types.Message {

	resp, err := getLPMessages(context.TODO(), queue.sqsClient, queue.mInput)

	if err != nil {
		fmt.Println("Got an error receiving messages:")
		fmt.Println(err)
		return nil
	}

	return resp.Messages
}

func (queue *queueSQS) RemoveMessage(msg types.Message) error {

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      queue.queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	}

	_, err := removeMessage(context.TODO(), queue.sqsClient, dMInput)

	if err != nil {
		err = errors.New("Got an error deleting the message from the queue: " + err.Error())
	}

	return err

}
