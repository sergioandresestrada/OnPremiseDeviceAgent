package service

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// ClientHBPort is an arbitrary port used in which the device API is listening
const ClientHBPort = "55555"

// Heartbeat receives a Message and prints it to stdout
// Returns a non-nil error if there's one during the execution and nil otherwise
func (s *Service) Heartbeat(msg Message) error {
	fmt.Println("Processing Heartbeat")
	if msg.Message == "" || msg.IPAddress == "" {
		err := errors.New("some message's expected fields are missing")
		return err
	}

	err := sendToClient(msg)
	return err
}

func sendToClient(message Message) error {
	client := net.ParseIP(message.IPAddress)
	if client == nil {
		return errors.New("invalid client IP")
	}
	host := "http://" + client.String()
	port := ClientHBPort

	fmt.Printf("sending heartbeat to %s.\n", host)

	res, err := http.Post(host+":"+port+"/heartbeat", "text/plain", bytes.NewBufferString(message.Message))

	if err != nil {
		err = fmt.Errorf("error performing the petition: %w", err)
		return err
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("error in the response: status code -> %v", res.StatusCode)
		return err
	}

	fmt.Println("Heartbeat sent and response received correctly.")
	return nil

}
