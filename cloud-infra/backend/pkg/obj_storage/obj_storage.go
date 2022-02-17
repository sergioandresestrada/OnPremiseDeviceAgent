package objstorage

import "mime/multipart"

type Obj_storage interface {
	UploadFile(*multipart.File, string) error
}
