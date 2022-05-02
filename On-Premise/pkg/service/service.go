package service

import (
	"On-Premise/pkg/config"
	objstorage "On-Premise/pkg/obj_storage"
	"On-Premise/pkg/queue"
	"On-Premise/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Message is just a reference to type Message in package types so that the usage is shorter
type Message = types.Message

// JobClient is just a reference to type JobClient in package types so that the usage is shorter
type JobClient = types.JobClient

// Service is the struct used to set up the On-Premise Server
// It contains a queue and object storage implementation
type Service struct {
	queue      queue.Queue
	objStorage objstorage.ObjStorage
}

// NewService creates and returns the reference to a new Service struct
func NewService(queue queue.Queue, objStorage objstorage.ObjStorage) *Service {
	s := &Service{
		queue:      queue,
		objStorage: objStorage,
	}
	return s
}

// Run is the main program loop.
// It will poll for messages from the queue and process them one by one
func (s *Service) Run() {
	for {
		receivedMessages := s.queue.ReceiveMessages()

		for _, queueMsg := range receivedMessages {
			var parsedMessage Message
			err := json.Unmarshal([]byte(*queueMsg.Body), &parsedMessage)
			if err != nil {
				fmt.Println("Error while unmarshalling the message")
				continue
			}

			go s.processMessage(parsedMessage)

			err = s.queue.RemoveMessage(queueMsg)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}

			fmt.Printf("Message was read and deleted successfully\n\n")
		}

	}
}

func (s *Service) processMessage(msg Message) {

	waitTime := config.InitialTimeBetweenRetries

	for i := 0; i < config.NumberOfRetries; i++ {
		var err error

		switch msg.Type {
		case "HEARTBEAT":
			err = s.Heartbeat(msg)
		case "JOB":
			err = s.Job(msg)
		case "UPLOAD":
			err = s.Upload(msg)
		default:
			fmt.Println("The received message is invalid")
			return
		}

		// if there was no error, we finished the processing, check for a url to send response and do it if present
		if err == nil {
			if msg.ResultURL != "" {
				s.sendMessageOutcome(msg, "SUCCESS")
			}
			break
		}

		// Otherwise, we log the error, send the result, wait the correspoding time and double it for next iteration
		fmt.Printf("There was an error processing the message: %v\n", err)

		s.sendMessageOutcome(msg, fmt.Sprintf("FAILURE: %v", err))

		time.Sleep(time.Duration(waitTime) * time.Second)
		waitTime *= 2
	}
}

func (s *Service) sendMessageOutcome(msg Message, result string) {

	url := msg.ResultURL + "/" + msg.DeviceUUID + "/" + msg.MessageUUID

	values := map[string]interface{}{
		"Result":    result,
		"Timestamp": time.Now().UnixMilli(),
	}

	jsonData, err := json.Marshal(values)

	if err != nil {
		fmt.Println("Error creating the result JSON to send")
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("There was an error sending the result or the server responded with status code different to 200")
	}
}
