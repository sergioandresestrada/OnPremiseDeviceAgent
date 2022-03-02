package objstorage

import (
	"io"
)

type Obj_storage interface {
	UploadFile(io.Reader, string) error
}
