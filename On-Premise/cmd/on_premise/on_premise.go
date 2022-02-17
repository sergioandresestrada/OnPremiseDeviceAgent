package main

import (
	"On-Premise/pkg/obj_storage"
	"On-Premise/pkg/queue"

	"On-Premise/pkg/service"
)

func setUpService() {
	queue := queue.NewQueueSQS()
	obj_storage := obj_storage.NewObjStorageS3()
	service := service.NewService(queue, obj_storage)
	service.Run()
}

func main() {
	setUpService()
}
