package api

import (
	"io/ioutil"
	"net/http"
)

func (s *Server) Jobs(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadFile("../../files/jobs.json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octec-stream")
	w.Write(files)
}

func (s *Server) Identification(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadFile("../../files/identification.json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(files)
}
