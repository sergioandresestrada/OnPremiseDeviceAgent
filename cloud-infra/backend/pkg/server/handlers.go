package server

import (
	"backend/pkg/types"
	"backend/pkg/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Message = types.Message

const BACKEND_URL = "http://192.168.1.208:12345"

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
	err = json.Unmarshal(requestBody, &message)
	if err != nil {
		fmt.Println("Invalid JSON provided as body")
		utils.BadRequest(w)
		return
	}

	fmt.Printf("\nrequestBody: %s\n", requestBody)
	fmt.Printf("Message content received: %v\n", message.Message)
	fmt.Printf("Type: %v\n", message.Type)

	if message.Message == "" || message.Type != "HEARTBEAT" {
		utils.BadRequest(w)
		return
	}

	err = s.queue.SendMessage(string(requestBody))
	if err != nil {
		fmt.Printf("%v\n", err)
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
	err = json.Unmarshal([]byte(r.FormValue("data")), &message)
	if err != nil {
		fmt.Println("Invalid JSON provided as data")
		utils.BadRequest(w)
		return
	}
	fmt.Printf("\nrequestBody: %s\n", r.FormValue("data"))
	fmt.Printf("Type: %v\n", message.Type)

	if message.Type != "JOB" || message.IPAddress == "" || message.Material == "" {
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateMaterial(message.Material)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateIPAddress(message.IPAddress)
	if err != nil {
		fmt.Printf("%v\n", err)
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
		fmt.Printf("%v\n", err)
		utils.BadRequest(w)
		return
	}

	rand.Seed(time.Now().UnixNano())
	message.FileName = fileHeader.Filename
	message.S3Name = strconv.Itoa(rand.Int())

	err = s.obj_storage.UploadFile(file, message.S3Name)

	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Got an error creating the message to the queue: %v\n", err)
		utils.ServerError(w)
		return
	}

	err = s.queue.SendMessage(string(messageJSON))
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)
}

func (s *Server) Upload(w http.ResponseWriter, r *http.Request) {
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

	if r.Header.Get("Content-Type") != "application/json" {
		fmt.Println("Invalid request content type")
		utils.BadRequest(w)
		return
	}

	var message Message
	err = json.Unmarshal(requestBody, &message)
	if err != nil {
		fmt.Println("Invalid JSON provided as body")
		utils.BadRequest(w)
		return
	}

	if message.Type != "UPLOAD" || message.IPAddress == "" || message.UploadInfo == "" {
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateIPAddress(message.IPAddress)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateUploadInfo(message.UploadInfo)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.BadRequest(w)
		return
	}

	fmt.Printf("\nReceived Type: %v\n", message.Type)
	fmt.Printf("Requested information: %v\n", message.UploadInfo)
	fmt.Printf("Device to request info from: %v\n", message.IPAddress)

	message.UploadURL = BACKEND_URL + "/upload" + message.UploadInfo

	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Got an error creating the message to the queue: %v\n", err)
		utils.ServerError(w)
		return
	}

	err = s.queue.SendMessage(string(messageJSON))
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)

}

func (s *Server) UploadIdentification(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		fmt.Println("Invalid request content type")
		utils.BadRequest(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading request body")
		utils.BadRequest(w)
		return
	}

	if !json.Valid(body) {
		fmt.Println("Invalid JSON as body")
		utils.BadRequest(w)
		return
	}

	deviceIP := r.Header.Get("X-Device")

	err = utils.ValidateIPAddress(deviceIP)

	if err != nil {
		fmt.Println("Device IP Header missing or invalid IP address in the request")
		utils.BadRequest(w)
		return
	}

	fmt.Printf("\nReceived Identification JSON from device: %v\n", deviceIP)

	fileName := "Identification-" + strings.Replace(deviceIP, ".", "_", 4) + ".json"
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error while creating the file")
		utils.BadRequest(w)
		return
	}

	defer file.Close()
	defer os.Remove(file.Name())

	io.Copy(file, bytes.NewBuffer(body))
	file.Seek(0, 0)

	err = s.obj_storage.UploadFile(file, fileName)

	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)

}

func (s *Server) UploadJobs(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		fmt.Println("Invalid request content type")
		utils.BadRequest(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading request body")
		utils.BadRequest(w)
		return
	}

	if !json.Valid(body) {
		fmt.Println("Invalid JSON as body")
		utils.BadRequest(w)
		return
	}

	deviceIP := r.Header.Get("X-Device")

	err = utils.ValidateIPAddress(deviceIP)

	if err != nil {
		fmt.Println("Device IP Header missing or invalid IP address in the request")
		utils.BadRequest(w)
		return
	}

	fmt.Printf("\nReceived Jobs JSON from device: %v\n", deviceIP)

	fileName := "Jobs-" + strings.Replace(deviceIP, ".", "_", 4) + ".json"
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error while creating the file")
		utils.BadRequest(w)
		return
	}

	defer file.Close()
	defer os.Remove(file.Name())

	io.Copy(file, bytes.NewBuffer(body))
	file.Seek(0, 0)

	err = s.obj_storage.UploadFile(file, fileName)

	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)
}
