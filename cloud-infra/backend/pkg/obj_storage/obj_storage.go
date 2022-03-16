package objstorage

import (
	"io"
)

// ObjStorage interface defines the methods that ObjStorage implementations will need to have
// Iterface is used although only one implementation is used so that we can mock it
type ObjStorage interface {
	UploadFile(io.Reader, string) error
}
