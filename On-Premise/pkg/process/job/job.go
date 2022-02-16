package job

import (
	"On-Premise/pkg/obj_storage"
	"On-Premise/pkg/types"
	"errors"
	"os"
)

type Message = types.Message

func ProcessJob(msg Message) error {

	fd, err := os.Create(msg.FileName)
	if err != nil {
		err = errors.New("Error while creating the file: " + err.Error())
		return err
	}
	defer fd.Close()

	err = obj_storage.DownloadFile(msg, fd)

	return err
}
