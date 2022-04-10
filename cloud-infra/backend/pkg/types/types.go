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
	DeviceName string `json:"DeviceName"`
}

// Information struct represents the names of the available files with information about the devices
// that are present in the object storage
type Information struct {
	Jobs           []string
	Identification []string
}

// Device struct representes the information about a device that we have, readed from the Database or received
// from an API call to create and store a new one
type Device struct {
	DeviceUUID string `json:"DeviceUUID,omitempty"`
	IP         string `json:"IP"`
	Name       string `json:"Name"`
	Model      string `json:"Model,omitempty"`
}
