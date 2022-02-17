package obj_storage

import (
	"On-Premise/pkg/types"
	"os"
)

type Message = types.Message

type Obj_storage interface {
	DownloadFile(Message, *os.File) error
}
