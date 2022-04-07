package main

import (
	"backend/pkg/database"
	objstorage "backend/pkg/obj_storage"
	"backend/pkg/queue"
	"backend/pkg/server"

	"github.com/gorilla/mux"
)

func setUpServer() {
	router := mux.NewRouter()
	queue := queue.NewQueueSQS()
	objstorage := objstorage.NewObjStorageS3()
	database := database.NewDatabaseDynamoDB()
	server := server.NewServer(queue, objstorage, database, router)

	server.Routes()
	server.ListenAndServe()

}

func main() {
	setUpServer()
}
