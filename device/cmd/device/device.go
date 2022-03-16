package main

import (
	"device/pkg/api"

	"github.com/gorilla/mux"
)

func setUpDevice() {
	router := mux.NewRouter()
	server := api.NewServer(router)

	server.Routes()
	server.ListenAndServe()
}

func main() {
	setUpDevice()
}
