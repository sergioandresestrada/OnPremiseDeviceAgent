package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
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

func GetQueueURL(c context.Context, api SQSGetLPMsgAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func GetLPMessages(c context.Context, api SQSGetLPMsgAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

func RemoveMessage(c context.Context, api SQSDeleteMessageAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}

func enviaCliente(message string) {
	host := "localhost"
	port := "9999"
	conType := "tcp"

	fmt.Printf("Connecting to %s on port %s.\n", host, port)

	conn, err := net.Dial(conType, host+":"+port)

	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Connection established correctly")

	_, err = conn.Write([]byte(message))

	if err != nil {
		fmt.Println("Error sending message:", err.Error())
		os.Exit(1)
	}

}

type Message struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	FileName string `json:"filename,omitempty"`
	S3Name   string `json:"s3name,omitempty"`
}

func main() {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	SQSclient := sqs.NewFromConfig(cfg)

	queueName := aws.String("test.fifo")
	qInput := &sqs.GetQueueUrlInput{
		QueueName: queueName,
	}

	result, err := GetQueueURL(context.TODO(), SQSclient, qInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	fmt.Printf("queueURL: %v\n", *queueURL)

	waitTime := 20

	for {

		mInput := &sqs.ReceiveMessageInput{
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

		resp, err := GetLPMessages(context.TODO(), SQSclient, mInput)
		if err != nil {
			fmt.Println("Got an error receiving messages:")
			fmt.Println(err)
			return
		}

		S3client := s3.NewFromConfig(cfg)

		for _, msg := range resp.Messages {
			fmt.Println("Message ID\t" + *msg.MessageId)
			var message Message
			json.Unmarshal([]byte(*msg.Body), &message)
			/* fmt.Printf("Type of message: %v\n", message.Type)
			fmt.Printf("Body of message: %v\n", message.Message)
			fmt.Printf("message.FileName: %v\n", message.FileName)
			fmt.Printf("message.S3Name: %v\n", message.S3Name) */

			switch message.Type {
			case "HEARTBEAT":
				processHeartbeatJob(message)
			case "FILE":
				processFileJob(message, S3client)
			}

			dMInput := &sqs.DeleteMessageInput{
				QueueUrl:      queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			}

			_, err = RemoveMessage(context.TODO(), SQSclient, dMInput)
			if err != nil {
				fmt.Println("Got an error deleting the message:")
				fmt.Println(err)
				return
			}

			fmt.Printf("Message was processed and deleted successfully\n\n")

		}
	}

}

func processHeartbeatJob(message Message) {
	fmt.Println("Processing Heartbeat Job")
	enviaCliente(message.Message)
}

func processFileJob(message Message, client *s3.Client) {
	fmt.Println("Processing file attached job")

	downloader := manager.NewDownloader(client)

	fd, err := os.Create(message.FileName)
	if err != nil {
		fmt.Println("Error while creating the file")
		return
	}
	defer fd.Close()

	fmt.Printf("Downloading file %s\n", message.FileName)

	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{
		Bucket: aws.String("sergiotfgbucket"),
		Key:    aws.String(message.S3Name),
	})

	if err != nil {
		fmt.Println("Error while downloading the file")
		return
	}

}
