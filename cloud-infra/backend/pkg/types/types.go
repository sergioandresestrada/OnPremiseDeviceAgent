package types

// Message struct represent the message with all its possible fields that any of the backend endpoints
// will probably receive
type Message struct {
	Type       string `json:"type"`
	Message    string `json:"message,omitempty"`
	FileName   string `json:"filename,omitempty"`
	S3Name     string `json:"s3name,omitempty"`
	Material   string `json:"material,omitempty"`
	IPAddress  string `json:"IPAddress,omitempty"`
	UploadInfo string `json:"UploadInfo,omitempty"`
	UploadURL  string `json:"UploadURL,omitempty"`
}

// Information struct represents the names of the available files with information about the devices
// that are present in the object storage
type Information struct {
	Jobs           []string
	Identification []string
}
