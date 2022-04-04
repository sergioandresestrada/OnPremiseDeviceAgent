package objstorage

import (
	"backend/pkg/types"
	"io"
	"os"
)

// ObjStorage interface defines the methods that ObjStorage implementations will need to have
// Iterface is used although only one implementation is used so that we can mock it
type ObjStorage interface {
	UploadFile(io.Reader, string) error
	AvailableInformation() (types.Information, error)
	GetFile(string, *os.File) error
}
