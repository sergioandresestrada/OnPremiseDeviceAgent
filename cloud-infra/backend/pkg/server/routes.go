package server

import (
	"fmt"
	"net/http"
)

func (s *Server) Routes() {
	s.router.HandleFunc("/", hello)
	s.router.HandleFunc("/heartbeat", s.Heartbeat).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/job", s.Job).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/upload", s.Upload).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/uploadIdentification", s.UploadIdentification).Methods("POST")
	s.router.HandleFunc("/uploadJobs", s.UploadJobs).Methods("POST")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "Working")
}
