package server

import (
	"backend/pkg/types"
	"backend/pkg/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Message = types.Message

func (s *Server) Heartbeat(w http.ResponseWriter, r *http.Request) {
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

	err = s.queue.SendMessage(string(requestBody))
	if err != nil {
		fmt.Println(err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)
}

func (s *Server) Job(w http.ResponseWriter, r *http.Request) {
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

	if message.Message == "" || message.Type != "JOB" || message.IPAddress == "" || message.Material == "" {
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateMaterial(message.Material)
	if err != nil {
		fmt.Println(err.Error())
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateIPAddress(message.IPAddress)
	if err != nil {
		fmt.Println(err.Error())
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

	err = utils.ValidateFile(file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println(err.Error())
		utils.BadRequest(w)
		return
	}

	rand.Seed(time.Now().UnixNano())
	message.FileName = fileHeader.Filename
	message.S3Name = strconv.Itoa(rand.Int())

	err = s.obj_storage.UploadFile(&file, message.S3Name)

	if err != nil {
		fmt.Println(err)
		utils.ServerError(w)
		return
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Got an error creating the message to the queue:")
		fmt.Println(err)
		utils.ServerError(w)
		return
	}

	err = s.queue.SendMessage(string(messageJSON))
	if err != nil {
		fmt.Println(err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)
}
