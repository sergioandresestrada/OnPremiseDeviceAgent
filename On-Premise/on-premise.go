package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
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

func GetQueueURL(c context.Context, api SQSGetLPMsgAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func GetLPMessages(c context.Context, api SQSGetLPMsgAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
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
	Type    string `json:"type"`
	Message string `json:"message"`
}

func main() {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)

	queueName := aws.String("test.fifo")
	qInput := &sqs.GetQueueUrlInput{
		QueueName: queueName,
	}

	result, err := GetQueueURL(context.TODO(), client, qInput)
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

		resp, err := GetLPMessages(context.TODO(), client, mInput)
		if err != nil {
			fmt.Println("Got an error receiving messages:")
			fmt.Println(err)
			return
		}

		fmt.Println("Message IDs:")

		for _, msg := range resp.Messages {
			fmt.Println("    " + *msg.MessageId)
			var message Message
			json.Unmarshal([]byte(*msg.Body), &message)
			fmt.Printf("Tipo del mensaje: %v\n", message.Type)
			fmt.Printf("Cuerpo del mensaje : %v\n", message.Message)
			enviaCliente(message.Message)
		}
	}

}
