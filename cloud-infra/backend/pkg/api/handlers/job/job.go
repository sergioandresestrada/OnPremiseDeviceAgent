package job

import (
	objstorage "backend/pkg/obj_storage"
	"backend/pkg/queue"
	"backend/pkg/types"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Message = types.Message

func Job(w http.ResponseWriter, r *http.Request) {
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

	if err != nil {
		fmt.Println("Error while reading the file")
	}

	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	message.FileName = fileHeader.Filename
	message.S3Name = strconv.Itoa(rand.Int())

	objstorage.UploadFile(&file, message.S3Name)

	s, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Got an error creating the message to the queue:")
		fmt.Println(err)
		return
	}
	queue.SendMessageToQueue(string(s))

}
