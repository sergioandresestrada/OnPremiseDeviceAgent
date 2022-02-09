package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Message struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	FileName string `json:"filename,omitempty"`
	S3Name   string `json:"s3name,omitempty"`
}

type SQSSendMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func GetQueueURL(c context.Context, api SQSSendMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func SendMsg(c context.Context, api SQSSendMessageAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(c, input)
}

func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func sendMessageToQueue(s string) {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)

	queueName := aws.String("test.fifo")
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: queueName,
	}

	result, err := GetQueueURL(context.TODO(), client, gQInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	sMInput := &sqs.SendMessageInput{

		MessageBody:    aws.String(s),
		QueueUrl:       queueURL,
		MessageGroupId: aws.String("1"),
	}

	resp, err := SendMsg(context.TODO(), client, sMInput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return
	}

	fmt.Println("Sent message with ID: " + *resp.MessageId)

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "Working")
}

func receiveMessage(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading request body")
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	var message Message
	json.Unmarshal(requestBody, &message)
	fmt.Printf("requestBody: %s\n", requestBody)
	fmt.Printf("Message content received: %v\n", message.Message)
	fmt.Printf("Type: %v\n", message.Type)

	sendMessageToQueue(string(requestBody))
}

func jobWithFile(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(64 << 20)

	if err != nil {
		fmt.Println("Error while reading request body")
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	var message Message
	json.Unmarshal([]byte(r.FormValue("data")), &message)
	fmt.Printf("requestBody: %s\n", r.FormValue("data"))
	fmt.Printf("Message content received: %v\n", message.Message)
	fmt.Printf("Type: %v\n", message.Type)

	file, fileHeader, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		fmt.Println("Error while reading the file")
	}

	BUCKETNAME := aws.String("sergiotfgbucket")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	rand.Seed(time.Now().UnixNano())
	message.FileName = fileHeader.Filename
	message.S3Name = strconv.Itoa(rand.Int())

	input := &s3.PutObjectInput{
		Bucket: BUCKETNAME,
		Key:    aws.String(message.S3Name),
		Body:   file,
	}

	_, err = PutFile(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error uploading file:")
		fmt.Println(err)
		return
	}

	s, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Got an error creating the message to the queue:")
		fmt.Println(err)
		return
	}
	sendMessageToQueue(string(s))

}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/", hello)
	router.HandleFunc("/message", receiveMessage).Methods("POST", "OPTIONS")
	router.HandleFunc("/jobwithfile", jobWithFile).Methods("POST", "OPTIONS")
	log.Fatal(http.ListenAndServe(":12345", router))
}

func main() {
	handleRequests()
}
