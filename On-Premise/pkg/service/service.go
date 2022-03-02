package service

import (
	"On-Premise/pkg/obj_storage"
	"On-Premise/pkg/queue"
	"On-Premise/pkg/types"
	"encoding/json"
	"fmt"
)

type Message = types.Message
type JobClient = types.JobClient

type Service struct {
	queue       queue.Queue
	obj_storage obj_storage.Obj_storage
}

func NewService(queue queue.Queue, obj_storage obj_storage.Obj_storage) *Service {
	s := &Service{
		queue:       queue,
		obj_storage: obj_storage,
	}
	return s
}

func (s *Service) Run() {
	for {
		receivedMessages := s.queue.ReceiveMessages()

		for _, queueMsg := range receivedMessages {
			var parsedMessage Message
			json.Unmarshal([]byte(*queueMsg.Body), &parsedMessage)

			var err error = nil

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
