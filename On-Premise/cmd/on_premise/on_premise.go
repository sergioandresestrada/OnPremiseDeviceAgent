package main

import (
	objstorage "On-Premise/pkg/obj_storage"
	"On-Premise/pkg/queue"

	"On-Premise/pkg/service"
)

func setUpService() {
	queue := queue.NewQueueSQS()
	objStorage := objstorage.NewObjStorageS3()
	service := service.NewService(queue, objStorage)
	service.Run()
}

func main() {
	setUpService()
}
