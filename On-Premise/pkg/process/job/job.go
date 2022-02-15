package job

import (
	"On-Premise/pkg/obj_storage"
	"On-Premise/pkg/types"
	"fmt"
	"os"
)

type Message = types.Message

func ProcessJob(msg Message) {

	fd, err := os.Create(msg.FileName)
	if err != nil {
		fmt.Println("Error while creating the file")
		return
	}
	defer fd.Close()

	obj_storage.DownloadFile(msg, fd)
}
