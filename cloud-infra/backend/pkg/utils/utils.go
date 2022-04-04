package utils

import (
	"errors"
	"io"
	"net"
	"net/http"
	"path/filepath"

	"github.com/hschendel/stl"
)

// BadRequest writes needed headers and status code 400 to the received http.ResponseWriter
func BadRequest(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusBadRequest)
}

// OKRequest writes needed headers and status code 200 to the received http.ResponseWriter
func OKRequest(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// ServerError writes needed headers and status code 500 to the received http.ResponseWriter
func ServerError(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusInternalServerError)
}

// ValidateFile checks whether the provided file is valid or not
// Returns nil if file is a valid .pdf or .stl file and an non-nil error otherwise
func ValidateFile(file io.ReadSeeker, FileName string, MIMEType string) error {

	switch filepath.Ext(FileName) {
	case ".pdf":
		if MIMEType == "application/pdf" {
			return nil
		}
	case ".stl":
		_, err := stl.ReadAll(file)
		if err != nil {
			return err
		}

		// need to set the reader back to the start of the file
		_, err = file.Seek(0, 0)
		return err
	}

	return errors.New("invalid file received")
}

// ValidateMaterial checks that the provided material is one of the valids one
// Returns nil if valid and a non-nil error otherwise
func ValidateMaterial(material string) error {
	switch material {
	case "HR PA 11", "HR PA 12GB", "HR TPA", "HR PP", "HR PA 12":
		return nil
	default:
		return errors.New("invalid material received")
	}
}

// ValidateIPAddress checks that the provided IP address is valid
// Returns nil if valid and a non-nil error otherwise
func ValidateIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip received")
	}
	return nil
}

// ValidateUploadInfo checks that the provided info type is valid
// Returns nil if valid and a non-nil error otherwise
func ValidateUploadInfo(info string) error {
	switch info {
	case "Jobs", "Identification":
		return nil
	default:
		return errors.New("invalid info requested")
	}
}
