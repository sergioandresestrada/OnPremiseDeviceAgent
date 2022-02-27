package service

import (
	"errors"
	"fmt"
	"net"
)

func (s *Service) Heartbeat(msg Message) error {
	fmt.Println("Processing Heartbeat Job")
	if msg.Message == "" {
		err := errors.New("some message's expected fields are missing")
		return err
	}

	err := sendToClient(msg.Message)
	return err
}

func sendToClient(message string) error {
	host := "localhost"
	port := "9999"
	conType := "tcp"

	fmt.Printf("Connecting to %s on port %s.\n", host, port)

	conn, err := net.Dial(conType, host+":"+port)

	if err != nil {
		err = errors.New("Error connecting:" + err.Error())
		return err
	}

	fmt.Println("Connection established correctly")

	_, err = conn.Write([]byte(message))

	if err != nil {
		err = errors.New("Error sending message:" + err.Error())
		return err
	}

	// At this point err will be nil always
	return err
}
