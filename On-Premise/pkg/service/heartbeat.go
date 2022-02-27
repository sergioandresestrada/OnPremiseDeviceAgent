package service

import (
	"errors"
	"fmt"
	"net"
)

const CLIENT_HB_PORT = "44444"

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
	port := CLIENT_HB_PORT
	conType := "tcp"

	fmt.Printf("Connecting to %s on port %s.\n", host, port)

	conn, err := net.Dial(conType, host+":"+port)

	if err != nil {
		err = fmt.Errorf("error connecting: %w", err)
		return err
	}

	defer conn.Close()

	fmt.Println("Connection established correctly")

	_, err = conn.Write([]byte(message))

	if err != nil {
		err = fmt.Errorf("error sending message: %w", err)
		return err
	}

	// At this point err will be nil always
	return err
}
