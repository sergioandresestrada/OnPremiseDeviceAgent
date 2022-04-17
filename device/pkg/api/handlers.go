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

// Jobs is the handler used with GETS /jobs endpoint
// It returns the information stored in ./files/jobs.json as body
// Will return status code 500 if there is a problem reading the mentioned file
func (s *Server) Jobs(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadFile("./files/jobs.json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(files)
	if err != nil {
		fmt.Println("error writing the Jobs file in the petition")
		return
	}

	fmt.Println("Served Jobs JSON file")
}

// Identification is the handler used with GETS /identification endpoint
// It returns the information stored in ./identification/jobs.json as body
// Will return status code 500 if there is a problem reading the mentioned file
func (s *Server) Identification(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadFile("./files/identification.json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(files)
	if err != nil {
		fmt.Println("error writing the Jobs file in the petition")
		return
	}

	fmt.Println("Served Identification JSON file")
}

// ReceiveJob is the handler used with POST /job endpoint
// It receives a Job as a MultipartForm request including JSON data in the 'job' field
// and a file in the 'file field'
// It will validate the received information and save the received file in the /receivedFiles folder
// It will return status code 200 or 400 as appropiate
func (s *Server) ReceiveJob(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(64 << 20)

	if err != nil {
		fmt.Println("Error while reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var job types.JobDevice
	err = json.Unmarshal([]byte(r.FormValue("job")), &job)
	if err != nil {
		fmt.Println("error while unmarshalling the received Job form")
	}

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

// Heartbeat is the handler used with POST /heartbeat endpoint
// It receives a single string as request body and prints it to stdout
// It will return status code 200 or 400 as appropiate
func (s *Server) Heartbeat(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Received Heartbeat: %s \n", requestBody)
}
