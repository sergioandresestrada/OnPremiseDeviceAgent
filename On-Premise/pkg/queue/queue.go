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

var sqsClient *sqs.Client
var queueURL *string
var mInput *sqs.ReceiveMessageInput

func getQueueURL(c context.Context, api SQSGetLPMsgAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func getLPMessages(c context.Context, api SQSGetLPMsgAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

func removeMessage(c context.Context, api SQSDeleteMessageAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}

func init() {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	sqsClient = sqs.NewFromConfig(cfg)

	queue := aws.String(queueName)

	qInput := &sqs.GetQueueUrlInput{
		QueueName: queue,
	}

	result, err := getQueueURL(context.TODO(), sqsClient, qInput)
	if err != nil {
		panic("Got an error getting the queue URL: " + err.Error())
	}

	queueURL = result.QueueUrl

	mInput = &sqs.ReceiveMessageInput{
		QueueUrl: queueURL,
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

func ReceiveMessages() []types.Message {

	resp, err := getLPMessages(context.TODO(), sqsClient, mInput)

	if err != nil {
		fmt.Println("Got an error receiving messages:")
		fmt.Println(err)
		return nil
	}

	return resp.Messages
}

func RemoveMessage(msg types.Message) error {

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	}

	_, err := removeMessage(context.TODO(), sqsClient, dMInput)

	if err != nil {
		err = errors.New("Got an error deleting the message from the queue: " + err.Error())
	}

	return err

}
