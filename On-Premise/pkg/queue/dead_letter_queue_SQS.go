package queue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
)

// DLQ_SQS defines the struct used to implement the dead letter queue interface using AWS SQS
// It contains an SQS client and the queue URL to be used
type DLQ_SQS struct {
	sqsClient *sqs.Client
	queueURL  *string
}

// NewDeadLetterQueueSQS creates and returns the reference to a new NewDeadLetterQueueSQS struct
func NewDeadLetterQueueSQS() *DLQ_SQS {
	q := &DLQ_SQS{}
	q.initialize()
	return q
}

func (dlq *DLQ_SQS) initialize() {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic(fmt.Sprintf("configuration error: %v", err))
	}

	dlq.sqsClient = sqs.NewFromConfig(cfg)

	queueNameString := aws.String(deadLetterQueueName)

	qInput := &sqs.GetQueueUrlInput{
		QueueName: queueNameString,
	}

	result, err := getQueueURL(context.TODO(), dlq.sqsClient, qInput)
	if err != nil {
		panic(fmt.Sprintf("Got an error getting the queue URL: %v", err))
	}

	dlq.queueURL = result.QueueUrl

}

// ReceiveMessages uses the queue to retrieve and return Messages from it
// Returns nil if there's an error receiving messages
func (dlq *DLQ_SQS) ReceiveMessages() []types.Message {

	mInput := &sqs.ReceiveMessageInput{
		QueueUrl: dlq.queueURL,
		AttributeNames: []types.QueueAttributeName{
			"SentTimestamp",
		},
		MaxNumberOfMessages: 10,
		MessageAttributeNames: []string{
			"All",
		},
		WaitTimeSeconds: int32(waitTime),
	}

	resp, err := getLPMessages(context.TODO(), dlq.sqsClient, mInput)

	if err != nil {
		fmt.Printf("Got an error receiving messages: %v\n", err)
		return nil
	}

	return resp.Messages
}

// RemoveMessage received a processed message and removes it from the queue
// Returns a non-nil error if there's one during the execution and nil otherwise
func (dlq *DLQ_SQS) RemoveMessage(msg types.Message) error {

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      dlq.queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	}

	_, err := removeMessage(context.TODO(), dlq.sqsClient, dMInput)

	if err != nil {
		err = fmt.Errorf("got an error deleting the message fron the queue: %w", err)
	}

	return err

}

// SendMessage receives an string and puts it in the correponding SQS URL
// Returns a non-nil error if there's one during the execution and nil otherwise
func (dlq *DLQ_SQS) SendMessage(s string) error {
	sMInput := &sqs.SendMessageInput{

		MessageBody:    aws.String(s),
		QueueUrl:       dlq.queueURL,
		MessageGroupId: aws.String("1"),
	}

	resp, err := sendMsg(context.TODO(), dlq.sqsClient, sMInput)
	if err != nil {
		err = fmt.Errorf("got an error sending the message to the Dead Letter Queue: %w", err)
		return err
	}

	fmt.Println("Sent message to Dead Letter Queue with ID: " + *resp.MessageId)
	return nil
}
