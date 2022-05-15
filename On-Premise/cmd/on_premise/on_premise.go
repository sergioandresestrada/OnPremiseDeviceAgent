package main

import (
	"On-Premise/pkg/config"
	objstorage "On-Premise/pkg/obj_storage"
	"On-Premise/pkg/queue"
	"On-Premise/pkg/types"
	"flag"

	"On-Premise/pkg/service"
)

func setUpService(config types.Config) {
	queue := queue.NewQueueSQS()
	objStorage := objstorage.NewObjStorageS3()
	service := service.NewService(queue, objStorage, config)
	service.Run()
}

func main() {

	numberRetries := flag.Int("r", config.NumberOfRetries, "The maximum number of retries when processing a message")
	secsBetweenRetries := flag.Int("s", config.InitialTimeBetweenRetries, "Time in seconds before the first retry (will double for successive retries)")

	flag.Parse()

	config := types.Config{
		NumberOfRetries:           *numberRetries,
		InitialTimeBetweenRetries: *secsBetweenRetries,
	}

	setUpService(config)
}
