package service

import (
	objstorage "On-Premise/pkg/obj_storage"
	"On-Premise/pkg/queue"
	"On-Premise/pkg/types"
	"encoding/json"
	"fmt"
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

			switch parsedMessage.Type {
			case "HEARTBEAT":
				err = s.Heartbeat(parsedMessage)
			case "JOB":
				err = s.Job(parsedMessage)
			case "UPLOAD":
				err = s.Upload(parsedMessage)
			default:
				fmt.Println("The received message is invalid")
				continue
			}

			if err != nil {
				fmt.Printf("There was an error processing the message: %v\n", err)
				continue
			}

			err = s.queue.RemoveMessage(queueMsg)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}

			fmt.Printf("Message was processed and deleted successfully\n\n")
		}

	}
}
