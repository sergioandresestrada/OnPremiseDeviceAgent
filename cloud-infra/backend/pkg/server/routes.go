package server

import (
	"fmt"
	"net/http"
)

// Routes defines the different endpoints the backend will have and assign the handlers to them
func (s *Server) Routes() {
	s.router.HandleFunc("/", hello)
	s.router.HandleFunc("/heartbeat", s.Heartbeat).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/job", s.Job).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/upload", s.Upload).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/uploadIdentification", s.UploadIdentification).Methods("POST")
	s.router.HandleFunc("/uploadJobs", s.UploadJobs).Methods("POST")
	s.router.HandleFunc("/availableInformation", s.AvailableInformation).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/getInformationFile", s.GetInformationFile).Methods("GET", "OPTIONS")

	// returns public info (name and model) from devices
	s.router.HandleFunc("/getPublicDevices", s.GetPublicDevices).Methods("GET", "OPTIONS")

	// CRUD funtionality for devices
	s.router.HandleFunc("/devices", s.DevicesCRUDOptionsHandler).Methods("OPTIONS")
	s.router.HandleFunc("/devices/{uuid}", s.DevicesCRUDOptionsHandler).Methods("OPTIONS")

	s.router.HandleFunc("/devices", s.GetDevices).Methods("GET")
	s.router.HandleFunc("/devices/{uuid}", s.GetDeviceByUUID).Methods("GET")
	s.router.HandleFunc("/devices", s.NewDevice).Methods("POST")
	s.router.HandleFunc("/devices/{uuid}", s.DeleteDevice).Methods("DELETE")
	s.router.HandleFunc("/devices/{uuid}", s.UpdateDevice).Methods("PUT")

	//Receives responses from the On Premise indicating the result of serving a message to the corresponding device
	s.router.HandleFunc("/responses/{deviceUUID}/{messageUUID}", s.ReceiveResponse).Methods("OPTIONS", "POST")

	// this are test handlers used to test UI without making unnecesary calls to AWS services
	s.router.HandleFunc("/testjobs", s.TestJobs).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/testidentification", s.TestIdentification).Methods("GET", "OPTIONS")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "Working")
}
