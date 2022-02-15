package heartbeat

import (
	"backend/pkg/queue"
	"backend/pkg/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Message = types.Message

func Heartbeat(w http.ResponseWriter, r *http.Request) {
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

	queue.SendMessageToQueue(string(requestBody))
}
