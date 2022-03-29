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
	s.router.HandleFunc("/testjobs", s.TestJobs).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/testidentification", s.TestIdentification).Methods("GET", "OPTIONS")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "Working")
}
