package utils

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/hschendel/stl"
)

func ValidateMaterial(material string) error {
	switch material {
	case "HR PA 11", "HR PA 12GB", "HR TPA", "HR PP", "HR PA 12":
		return nil
	default:
		return errors.New("invalid material received")
	}
}

func ValidateFile(file io.ReadSeeker, FileName string, MIMEType string) error {

	switch filepath.Ext(FileName) {
	case ".pdf":
		if MIMEType == "application/pdf" {
			return nil
		}
	case ".stl":
		_, err := stl.ReadAll(file)
		// need to set the reader back to the start of the file
		file.Seek(0, 0)
		return err
	}

	return errors.New("invalid file received")
}
