package main

import (
	"On-Premise/pkg/config"
	objstorage "On-Premise/pkg/obj_storage"
	"On-Premise/pkg/queue"
	"On-Premise/pkg/types"
	"flag"
	"fmt"

	"On-Premise/pkg/service"
)

func setUpService(config types.Config) {
	fmt.Println("Setting up...")
	messageQueue := queue.NewQueueSQS()
	objStorage := objstorage.NewObjStorageS3()
	DLQ := queue.NewDeadLetterQueueSQS()
	service := service.NewService(messageQueue, objStorage, DLQ, config)
	fmt.Println("Running correctly")
	service.Run()
}

func setUpDeadLetterQueueService() {
	fmt.Println("Setting up...")
	DLQ := queue.NewDeadLetterQueueSQS()
	service := service.NewDLQService(DLQ)
	fmt.Println("Running correctly")
	service.Run()
}

func main() {

	numberRetries := flag.Int("r", config.NumberOfRetries, "The maximum number of retries when processing a message")
	secsBetweenRetries := flag.Int("s", config.InitialTimeBetweenRetries, "Time in seconds before the first retry (will double for successive retries)")
	dlq := flag.Bool("dlq", false, "If set, reads, shows and deletes messages from the Dead Letter Queue")

	flag.Parse()

	if *dlq {
		setUpDeadLetterQueueService()
	}

	config := types.Config{
		NumberOfRetries:           *numberRetries,
		InitialTimeBetweenRetries: *secsBetweenRetries,
	}

	setUpService(config)
}
