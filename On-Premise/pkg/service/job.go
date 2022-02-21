package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const CLIENT_JOB_PORT = "55555"

func (s *Service) Job(msg Message) error {

	if msg.FileName == "" || msg.Message == "" || msg.S3Name == "" || msg.Material == "" || msg.IPAddress == "" {
		err := errors.New("some message's expected fields are missing")
		return err
	}

	fd, err := os.Create("onPremiseFiles/" + msg.FileName)
	if err != nil {
		err = errors.New("Error while creating the file: " + err.Error())
		return err
	}
	defer fd.Close()

	err = s.obj_storage.DownloadFile(msg, fd)
	if err != nil {
		err = errors.New("Error downloading the file: " + err.Error())
		return err
	}

	jobToClient := JobClient{}
	jobToClient.FileName = msg.FileName
	jobToClient.Material = msg.Material

	err = sendJobToClient(jobToClient, fd, msg.IPAddress)

	return err
}

func sendJobToClient(job JobClient, fd *os.File, clientIP string) error {
	client := net.ParseIP(clientIP)
	if client == nil {
		return errors.New("invalid client IP")
	}

	JobJson, err := json.Marshal(&job)

	if err != nil {
		return errors.New("error creating the job to sent to the client")
	}

	conn, err := net.Dial("tcp", client.String()+":"+CLIENT_JOB_PORT)
	if err != nil {
		return errors.New("error connecting to the device")
	}

	n, err := conn.Write(JobJson)
	fmt.Printf("sent %v bytes\n", n)
	if err != nil {
		return errors.New("error while sending the job")
	}

	// receive a single byte buffer for sync
	buffer := make([]byte, 1)
	conn.Read(buffer)

	_, err = io.Copy(conn, fd)
	if err != nil {
		return errors.New("error sending the file")
	}
	conn.Close()

	return nil
}
