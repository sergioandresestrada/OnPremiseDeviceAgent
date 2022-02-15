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

			switch message.Type {
			case "HEARTBEAT":
				hb.ProcessHeartbeat(message)
			case "FILE":
				job.ProcessJob(message)
			}
			queue.RemoveMessage(msg)
			fmt.Printf("Message was processed and deleted successfully\n\n")
		}

	}
}
