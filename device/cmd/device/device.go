package main

import (
	"device/pkg/api"
	"fmt"

	"github.com/gorilla/mux"
)

func setUpDevice() {
	fmt.Println("Setting up...")
	router := mux.NewRouter()
	server := api.NewServer(router)

	server.Routes()
	fmt.Println("Running correctly")
	server.ListenAndServe()
}

func main() {
	setUpDevice()
}
