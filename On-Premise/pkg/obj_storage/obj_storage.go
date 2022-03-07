package objstorage

import (
	"On-Premise/pkg/types"
	"os"
)

// Message is just a reference to type Message in package types so that the usage is shorter
type Message = types.Message

// ObjStorage interface defines the methods that ObjStorage implementations will need to have
// Iterface is used although only one implementation is used so that we can mock it
type ObjStorage interface {
	DownloadFile(Message, *os.File) error
}
