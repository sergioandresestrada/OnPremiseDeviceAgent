package main

import (
	"fmt"
	"log"
	"net/http"

	hb "backend/pkg/api/handlers/heartbeat"
	"backend/pkg/api/handlers/job"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "Working")
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/", hello)
	router.HandleFunc("/heartbeat", hb.Heartbeat).Methods("POST", "OPTIONS")
	router.HandleFunc("/job", job.Job).Methods("POST", "OPTIONS")
	log.Fatal(http.ListenAndServe(":12345", router))
}

func main() {
	handleRequests()
}
