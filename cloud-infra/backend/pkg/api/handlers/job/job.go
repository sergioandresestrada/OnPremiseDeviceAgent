package job

import (
	"backend/pkg/api/utils"
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

	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	var message Message
	json.Unmarshal([]byte(r.FormValue("data")), &message)
	fmt.Printf("requestBody: %s\n", r.FormValue("data"))
	fmt.Printf("Message content received: %v\n", message.Message)
	fmt.Printf("Type: %v\n", message.Type)

	if message.Message == "" || message.Type != "JOB" {
		utils.BadRequest(w)
		return
	}

	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		fmt.Println("Error while reading the file")
		utils.BadRequest(w)
		return
	}

	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	message.FileName = fileHeader.Filename
	message.S3Name = strconv.Itoa(rand.Int())

	err = objstorage.UploadFile(&file, message.S3Name)

	if err != nil {
		fmt.Println(err)
		utils.ServerError(w)
		return
	}

	s, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Got an error creating the message to the queue:")
		fmt.Println(err)
		utils.ServerError(w)
		return
	}

	err = queue.SendMessageToQueue(string(s))
	if err != nil {
		fmt.Println(err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)
}
