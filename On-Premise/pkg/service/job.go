package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ClientJobPort is an arbitrary port used in which the device API is listening
const ClientJobPort = "55555"

// Job receives a message, validate it fields and send it to the device using its API
// Returns a non-nil error if there's one during the execution and nil otherwise
func (s *Service) Job(msg Message) error {

	if msg.FileName == "" || msg.S3Name == "" || msg.Material == "" || msg.IPAddress == "" {
		err := errors.New("some message's expected fields are missing")
		return err
	}

	fd, err := os.CreateTemp("onPremiseFiles/", "")
	if err != nil {
		err = fmt.Errorf("error while creating the file: %w", err)
		return err
	}

	defer os.Remove(fd.Name())

	err = s.objStorage.DownloadFile(msg, fd)
	if err != nil {
		err = fmt.Errorf("error downloading the file: %w", err)
		return err
	}

	jobToClient := JobClient{}
	jobToClient.FileName = msg.FileName
	jobToClient.Material = msg.Material

	err = sendJobToClient(jobToClient, fd, msg.FileName, msg.IPAddress)

	return err
}

func sendJobToClient(job JobClient, fd *os.File, fileName string, clientIP string) error {
	client := net.ParseIP(clientIP)
	if client == nil {
		return errors.New("invalid client IP")
	}

	JobJSON, err := json.Marshal(&job)

	if err != nil {
		return errors.New("error creating the job to send to the client")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormField("job")
	if err != nil {
		return errors.New("error including the JSON in the petition")
	}

	_, err = io.Copy(fw, strings.NewReader(string(JobJSON)))
	if err != nil {
		return errors.New("error writing the JSON in the petition")
	}

	switch filepath.Ext(fileName) {
	case ".pdf":
		fw, err = customCreateFormFile(writer, "file", fileName, "application/pdf")
	default:
		fw, err = writer.CreateFormFile("file", fileName)
	}

	if err != nil {
		return errors.New("error including the file in the petition")
	}

	_, err = io.Copy(fw, fd)
	if err != nil {
		return errors.New("error writing the file in the petition")
	}

	writer.Close()

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", "http://"+clientIP+":"+ClientJobPort+"/job", body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, err := httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("Error while performing the request %v", err)
	}

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("resquest failed with status code %v", rsp.StatusCode)
	}

	return nil
}

func customCreateFormFile(w *multipart.Writer, fieldName string, fileName string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}
