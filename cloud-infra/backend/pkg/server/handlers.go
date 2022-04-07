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

// Message is just a reference to type Message in package types so that the usage is shorter
type Message = types.Message

// Device is just a reference to type Device in package types so that the usage is shorter
type Device = types.Device

// Heartbeat is the handler used with POST and OPTIONS /heartbeat endpoint
// It will validate the received JSON, if valid, and send the corresponding message to the queue
// It will return status code 200, 400 or 500 as appropiate
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

// Job is the handler used with POST and OPTIONS /job endpoint
// It will validate the received MultiPart Form, if valid,
// and send the corresponding message to the queue and file to object storage
// It will return status code 200, 400 or 500 as appropiate
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

	err = s.objStorage.UploadFile(file, message.S3Name)

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

// Upload is the handler used with POST and OPTIONS /upload endpoint
// It will validate the received JSON, if valid, and send the corresponding message to the queue,
// including the URL that the On-Premise server will have to use to upload the requested information
// It will return status code 200, 400 or 500 as appropiate
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

	message.UploadURL = s.serverURL + "/upload" + message.UploadInfo

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

// UploadIdentification is the handler used with POST /uploadIdentification endpoint
// It will receive a JSON body containing device's identification information
// and device's IP in the X-Device header, create the correspoding file
// and upload it to the object storage
// It will return status code 200, 400 or 500 as appropiate
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

	_, err = io.Copy(file, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	err = s.objStorage.UploadFile(file, fileName)

	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)

}

// UploadJobs is the handler used with POST /uploadJobs endpoint
// It will receive a JSON body containing device's jobs information
// and device's IP in the X-Device header, create the correspoding file
// and upload it to the object storage
// It will return status code 200, 400 or 500 as appropiate
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

	_, err = io.Copy(file, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	err = s.objStorage.UploadFile(file, fileName)

	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	utils.OKRequest(w)
}

// AvailableInformation is the handler used with GET /availableInformation endpoint
// It will return a JSON with all the Jobs and Identification information files
// that are available in the object storage
// It will return status code 200, 400 or 500 as appropiate
func (s *Server) AvailableInformation(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	AvailableInformation, err := s.objStorage.AvailableInformation()
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	jsonResult, err := json.Marshal(AvailableInformation)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(jsonResult)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}
	fmt.Printf("\nServed the list of Available Information\n")

}

// GetInformationFile is the handler used with GET /getInformationFile endpoint
// It will return the requestes JSON file, if it a valid one and it exists in the object storage
// Requested file name is received from the petition as a Get parameter
// It will return status code 200, 400 or 500 as appropiate
func (s *Server) GetInformationFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	key := r.URL.Query().Get("file")
	if !strings.HasPrefix(key, "Jobs-") && !strings.HasPrefix(key, "Identification-") || !strings.HasSuffix(key, ".json") {
		utils.BadRequest(w)
		return
	}

	file, err := os.CreateTemp("/tmp", "file")
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}
	defer os.Remove(file.Name())

	err = s.objStorage.GetFile(key, file)

	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}
	fmt.Printf("\nServed file %s\n", key)

}

func (s *Server) GetDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	devices, err := s.database.GetDevices()
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	publicJSON := utils.DevicesToPublicJSON(devices)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(publicJSON)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}
	fmt.Printf("\nServed the information of available Devices\n")

}

// TestJobs is the test handler used with GET /testjobs endpoint
func (s *Server) TestJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}
	file, _ := os.ReadFile("jobs.json")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err := w.Write(file)
	if err != nil {
		fmt.Printf("There was an error writing the information: %v\n", err)
		return
	}
}

// TestIdentification is the test handler used with GET /testidentification endpoint
func (s *Server) TestIdentification(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}
	file, _ := os.ReadFile("identification.json")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err := w.Write(file)
	if err != nil {
		fmt.Printf("There was an error writing the information: %v\n", err)
		return
	}
}
