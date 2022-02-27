package queue

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	queueName = "test.fifo"
)

type queueSQS struct {
	sqsClient *sqs.Client
	queueURL  *string
}

func NewQueueSQS() *queueSQS {
	q := &queueSQS{}
	q.initialize()
	return q
}

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

func (queue *queueSQS) initialize() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	queue.sqsClient = sqs.NewFromConfig(cfg)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	result, err := getQueueURL(context.TODO(), queue.sqsClient, gQInput)
	if err != nil {
		panic("Got an error getting the queue URL: " + err.Error())
	}

	queue.queueURL = result.QueueUrl
}

func (queue *queueSQS) SendMessage(s string) error {
	sMInput := &sqs.SendMessageInput{

		MessageBody:    aws.String(s),
		QueueUrl:       queue.queueURL,
		MessageGroupId: aws.String("1"),
	}

	resp, err := sendMsg(context.TODO(), queue.sqsClient, sMInput)
	if err != nil {
		err = errors.New("Got an error sending the message to the queue: " + err.Error())
		return err
	}

	fmt.Println("Sent message with ID: " + *resp.MessageId)
	return nil
}
