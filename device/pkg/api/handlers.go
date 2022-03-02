package api

import (
	"device/pkg/types"
	"device/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func (s *Server) Jobs(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadFile("./files/jobs.json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(files)

	fmt.Println("Served Jobs JSON file")
}

func (s *Server) Identification(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadFile("./files/identification.json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(files)

	fmt.Println("Served Identification JSON file")
}

func (s *Server) ReceiveJob(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(64 << 20)

	if err != nil {
		fmt.Println("Error while reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var job types.JobDevice
	json.Unmarshal([]byte(r.FormValue("job")), &job)

	if job.FileName == "" || job.Material == "" {
		fmt.Println("Missing field")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = utils.ValidateMaterial(job.Material)
	if err != nil {
		fmt.Printf("%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		fmt.Println("Error while reading the file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer file.Close()

	err = utils.ValidateFile(file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		fmt.Printf("%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	localFile, err := os.Create("./receivedFiles/" + job.FileName)

	if err != nil {
		fmt.Printf("Error while creating the file: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer localFile.Close()

	_, err = io.Copy(localFile, file)

	if err != nil {
		fmt.Printf("Error while saving the file: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Received new job with file %v and material %v\n", job.FileName, job.Material)
}

func (s *Server) Heartbeat(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Received Heartbeat: %s \n", requestBody)
}
