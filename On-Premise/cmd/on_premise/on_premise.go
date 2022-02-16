package main

import (
	hb "On-Premise/pkg/process/heartbeat"
	"On-Premise/pkg/process/job"
	"On-Premise/pkg/queue"
	"On-Premise/pkg/types"
	"encoding/json"
	"fmt"
)

type Message = types.Message

func main() {

	for {
		receivedMessages := queue.ReceiveMessages()

		for _, msg := range receivedMessages {
			var message Message
			json.Unmarshal([]byte(*msg.Body), &message)

			var err error = nil

			switch message.Type {
			case "HEARTBEAT":
				err = hb.ProcessHeartbeat(message)
			case "FILE":
				err = job.ProcessJob(message)
			}
			if err != nil {
				fmt.Printf("There was an error processing the message: %v\n", err.Error())
				continue
			}

			err = queue.RemoveMessage(msg)
			if err != nil {
				fmt.Printf(err.Error())
				continue
			}

			fmt.Printf("Message was processed and deleted successfully\n\n")
		}

	}
}
