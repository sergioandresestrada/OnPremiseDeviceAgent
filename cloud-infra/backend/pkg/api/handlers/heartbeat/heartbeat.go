package heartbeat

import (
	"backend/pkg/api/utils"
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
		utils.BadRequest(w)
		return
	}

	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	var message Message
	json.Unmarshal(requestBody, &message)
	fmt.Printf("requestBody: %s\n", requestBody)
	fmt.Printf("Message content received: %v\n", message.Message)
	fmt.Printf("Type: %v\n", message.Type)

	if message.Message == "" || message.Type != "HEARTBEAT" {
		utils.BadRequest(w)
		return
	}

	err = queue.SendMessageToQueue(string(requestBody))
	if err != nil {
		fmt.Println(err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)
}
