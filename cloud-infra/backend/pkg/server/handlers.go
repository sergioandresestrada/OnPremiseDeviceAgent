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

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	if message.DeviceName == "" {
		fmt.Println("Device name field missing")
		utils.BadRequest(w)
		return
	}

	fmt.Printf("\nrequestBody: %s\n", requestBody)
	fmt.Printf("Message content received: %v\n", message.Message)
	fmt.Printf("Type: %v\n", message.Type)
	fmt.Printf("Device name: %v\n", message.DeviceName)

	if message.Message == "" || message.Type != "HEARTBEAT" {
		utils.BadRequest(w)
		return
	}

	deviceIP, err := s.database.DeviceFromName(message.DeviceName)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	if deviceIP == "" {
		fmt.Printf("No device found with provided name\n")
		utils.BadRequest(w)
		return
	}

	message.IPAddress = deviceIP

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

	if message.Type != "JOB" || message.DeviceName == "" || message.Material == "" {
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateMaterial(message.Material)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.BadRequest(w)
		return
	}

	deviceIP, err := s.database.DeviceFromName(message.DeviceName)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	if deviceIP == "" {
		fmt.Printf("No device found with provided name\n")
		utils.BadRequest(w)
		return
	}

	message.IPAddress = deviceIP

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
	message.S3Name = strconv.Itoa(rand.Int()) + " - " + message.FileName

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

	if message.Type != "UPLOAD" || message.DeviceName == "" || message.UploadInfo == "" {
		utils.BadRequest(w)
		return
	}

	deviceIP, err := s.database.DeviceFromName(message.DeviceName)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	if deviceIP == "" {
		fmt.Printf("No device found with provided name\n")
		utils.BadRequest(w)
		return
	}

	message.IPAddress = deviceIP

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

	deviceName := r.Header.Get("X-Device")

	if deviceName == "" {
		fmt.Println("Device Name Header missing")
		utils.BadRequest(w)
		return
	}

	fmt.Printf("\nReceived Identification JSON from device: %v\n", deviceName)

	fileName := "Identification-" + strings.Replace(deviceName, ".", "_", 4) + ".json"
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

	deviceName := r.Header.Get("X-Device")

	if deviceName == "" {
		fmt.Println("Device Name Header missing")
		utils.BadRequest(w)
		return
	}

	fmt.Printf("\nReceived Jobs JSON from device: %v\n", deviceName)

	fileName := "Jobs-" + strings.Replace(deviceName, ".", "_", 4) + ".json"
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

// GetPublicDevices is the handler used with GET /getPublicDevices endpoint
// It will return the information about all the devices in JSON format
// It will return status code 200 or 500 as appropiate
func (s *Server) GetPublicDevices(w http.ResponseWriter, r *http.Request) {
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

	devicesJSON, err := json.Marshal(devices)
	if err != nil {
		fmt.Printf("Error while creating the JSON%v\n", err)
		utils.ServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(devicesJSON)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}
	fmt.Printf("\nServed the list of Devices\n")
}

func (s *Server) GetDeviceByUUID(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	deviceUUID := mux.Vars(r)["uuid"]

	if deviceUUID == "" {
		fmt.Printf("Missing device UUID in request\n")
		utils.BadRequest(w)
		return
	}

	device, err := s.database.GetDeviceByUUID(deviceUUID)

	if err != nil {
		fmt.Printf("Error while getting the device: %v\n", err)
		utils.ServerError(w)
		return
	}

	// Checks if the item was found or returned value is empty
	if device.Name == "" {
		fmt.Printf("Device not found with given UUID\n")
		utils.BadRequest(w)
		return
	}

	deviceJSON, err := json.Marshal(device)
	if err != nil {
		fmt.Printf("Error while creating the JSON%v\n", err)
		utils.ServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(deviceJSON)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}
	fmt.Printf("\nServed the information of device with UUID: %v\n", deviceUUID)

}

func (s *Server) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	deviceUUID := mux.Vars(r)["uuid"]

	if deviceUUID == "" {
		fmt.Printf("Missing device UUID in request\n")
		utils.BadRequest(w)
		return
	}

	err := s.database.DeleteDeviceFromUUID(deviceUUID)
	if err != nil {
		fmt.Printf("Error while deleting the device\n")
		utils.ServerError(w)
		return
	}

	fmt.Printf("Deleted device with UUID %v\n", deviceUUID)
	utils.OKRequest(w)
}

func (s *Server) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	deviceUUID := mux.Vars(r)["uuid"]

	if deviceUUID == "" {
		fmt.Printf("Missing device UUID in request\n")
		utils.BadRequest(w)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading request body")
		utils.BadRequest(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		fmt.Println("Invalid request content type")
		utils.BadRequest(w)
		return
	}

	var device Device
	err = json.Unmarshal(requestBody, &device)

	if err != nil {
		fmt.Println("New Device: Invalid JSON provided as body")
		utils.BadRequest(w)
		return
	}

	if device.IP == "" || device.Name == "" {
		fmt.Println("New Device: Invalid JSON provided as body, missing fields")
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateIPAddress(device.IP)
	if err != nil {
		fmt.Println("Invalid IP address received")
		utils.BadRequest(w)
		return
	}

	device.DeviceUUID = deviceUUID

	err = s.database.UpdateDevice(device)
	if err != nil {
		fmt.Printf("Error while updating: %v", err)
		utils.ServerError(w)
		return
	}

	fmt.Printf("Updated device with UUID %v\n", deviceUUID)
	utils.OKRequest(w)

}

// NewDevice is the handler used with POST /devices endpoint
// It preforms all the necessary checking and, if everything is correct, will insert a new device to the DB
// It will return status code 200, 400 or 500 as appropiate
func (s *Server) NewDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.OKRequest(w)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading request body")
		utils.BadRequest(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		fmt.Println("Invalid request content type")
		utils.BadRequest(w)
		return
	}

	var device Device
	err = json.Unmarshal(requestBody, &device)

	if err != nil {
		fmt.Println("New Device: Invalid JSON provided as body")
		utils.BadRequest(w)
		return
	}

	if device.IP == "" || device.Name == "" {
		fmt.Println("New Device: Invalid JSON provided as body, missing fields")
		utils.BadRequest(w)
		return
	}

	err = utils.ValidateIPAddress(device.IP)
	if err != nil {
		fmt.Println("Invalid IP address received")
		utils.BadRequest(w)
		return
	}

	exists, err := s.database.DeviceExistWithNameAndIP(device.Name, device.IP)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}

	if exists {
		fmt.Println("The device Name or IP provided already exist")
		utils.BadRequest(w)
		return
	}

	device.DeviceUUID = uuid.NewString()

	err = s.database.InsertDevice(device)
	if err != nil {
		fmt.Printf("%v\n", err)
		utils.ServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("Device inserted successfully")

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
