package service

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const CLIENT_PORT = "55555"

func (s *Service) Upload(msg Message) error {
	fmt.Println("Processing Upload")

	if msg.IPAddress == "" || msg.UploadInfo == "" || msg.UploadURL == "" {
		err := errors.New("some message's expected fields are missing")
		return err
	}

	buffer, err := receiveInfoFromDevice(msg)

	if err != nil {
		return fmt.Errorf("error receiving information from the device: %w", err)
	}

	err = sendInfoToBackend(buffer, msg.UploadURL, msg.IPAddress)
	if err != nil {
		return fmt.Errorf("error while sending the information to the backend: %w", err)
	}

	fmt.Println("Information obtained and sent correctly")

	return nil
}

func receiveInfoFromDevice(msg Message) ([]byte, error) {
	client := net.ParseIP(msg.IPAddress)
	if client == nil {
		return nil, errors.New("invalid client IP")
	}

	res, err := http.Get("http://" + client.String() + ":" + CLIENT_PORT + "/" + strings.ToLower(msg.UploadInfo))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected status code 200, got %v instead", res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("expected JSON body, got %v instead", res.Header.Get("Content-Type"))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error while reading response's body")
	}

	return body, nil
}

func sendInfoToBackend(info []byte, url string, deviceIP string) error {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(info))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Device", deviceIP)
	res, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("expected status code 200, got %v instead", res.StatusCode)
	}

	return nil
}
