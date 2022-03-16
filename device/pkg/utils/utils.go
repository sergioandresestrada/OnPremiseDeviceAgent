package utils

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/hschendel/stl"
)

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
		// need to set the reader back to the start of the file
		file.Seek(0, 0)
		return err
	}

	return errors.New("invalid file received")
}
