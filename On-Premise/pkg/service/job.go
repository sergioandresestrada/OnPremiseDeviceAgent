package service

import (
	"errors"
	"os"
)

func (s *Service) Job(msg Message) error {

	if msg.FileName == "" || msg.Message == "" || msg.S3Name == "" {
		err := errors.New("some message's expected fields are missing")
		return err
	}

	fd, err := os.Create(msg.FileName)
	if err != nil {
		err = errors.New("Error while creating the file: " + err.Error())
		return err
	}
	defer fd.Close()

	err = s.obj_storage.DownloadFile(msg, fd)

	return err
}
