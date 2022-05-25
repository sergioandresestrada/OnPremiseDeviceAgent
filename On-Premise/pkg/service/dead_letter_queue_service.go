package service

import (
	"On-Premise/pkg/queue"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type DLQService struct {
	queue queue.DeadLetterQueue
}

// NewDLQService creates and returns the reference to a new DLQService struct
func NewDLQService(queue queue.DeadLetterQueue) *DLQService {
	s := &DLQService{
		queue: queue,
	}
	return s
}

// Run is the main program loop.
// It will poll for messages from the Dead Letter Queue, read, show and delete them until there are no left messages in the queue.
func (s *DLQService) Run() {
	receivedMessages := s.queue.ReceiveMessages()
	for len(receivedMessages) > 0 {
		for _, queueMsg := range receivedMessages {
			var parsedMessage DLQ_Message
			err := json.Unmarshal([]byte(*queueMsg.Body), &parsedMessage)
			if err != nil {
				fmt.Println("Error while unmarshalling the message")
				continue
			}

			go s.showMessage(parsedMessage)

			err = s.queue.RemoveMessage(queueMsg)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}

			fmt.Printf("Message was deleted successfully\n\n\n")
		}
		receivedMessages = s.queue.ReceiveMessages()
	}

	fmt.Println("All messages have been shown and deleted. Dead Letter Queue is now empty.")
	fmt.Println("Exiting....")
	os.Exit(0)
}

func (s *DLQService) showMessage(msg DLQ_Message) {
	date := time.Unix(msg.Timestamp/1000, 0).In(time.Local).Format("02/01/2006 15:04:05")

	fmt.Printf("Message processed on %v\n", date)
	fmt.Printf("\tDevice Name: %v\n", msg.DeviceName)
	fmt.Printf("\tType: %v\n", msg.Type)
	switch msg.Type {
	case "HEARTBEAT":
		fmt.Printf("\tAssociated message: %v\n", msg.AdditionalInfo)
	case "JOB":
		fmt.Printf("\tAttached file: %v\n", msg.AdditionalInfo)
	case "UPLOAD":
		fmt.Printf("\t Requested information: %v\n", msg.AdditionalInfo)
	}
	fmt.Printf("\tLast Recorded Result: %v\n\n", msg.LastResult)
}