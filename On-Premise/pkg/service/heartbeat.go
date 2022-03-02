package service

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

const CLIENT_HB_PORT = "55555"

func (s *Service) Heartbeat(msg Message) error {
	fmt.Println("Processing Heartbeat")
	if msg.Message == "" {
		err := errors.New("some message's expected fields are missing")
		return err
	}

	err := sendToClient(msg.Message)
	return err
}

func sendToClient(message string) error {
	host := "http://127.0.0.1"
	port := CLIENT_HB_PORT

	fmt.Printf("sending heartbeat to %s.\n", host)

	res, err := http.Post(host+":"+port+"/heartbeat", "text/plain", bytes.NewBufferString(message))

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
