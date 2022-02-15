package types

type Message struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	FileName string `json:"filename,omitempty"`
	S3Name   string `json:"s3name,omitempty"`
}
