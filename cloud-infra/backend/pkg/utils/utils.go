package utils

import (
	"errors"
	"io"
	"net"
	"net/http"
	"path/filepath"

	"github.com/hschendel/stl"
)

func BadRequest(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusBadRequest)
}

func OKRequest(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func ServerError(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusInternalServerError)
}

func ValidateFile(file io.ReadSeeker, FileName string, MIMEType string) error {

	switch filepath.Ext(FileName) {
	case ".pdf":
		if MIMEType == "application/pdf" {
			return nil
		}
	case ".stl":
		_, err := stl.ReadAll(file)
		return err
	}

	return errors.New("invalid file received")
}

func ValidateMaterial(material string) error {
	switch material {
	case "HR PA 11", "HR PA 12GB", "HR TPA", "HR PP", "HR PA 12":
		return nil
	default:
		return errors.New("invalid material received")
	}
}

func ValidateIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip received")
	}
	return nil
}
