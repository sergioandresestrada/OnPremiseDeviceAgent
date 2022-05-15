package types

// Message struct represent the message with all its possible fields that any of the backend endpoints
// will probably receive
type Message struct {
	Type        string `json:"type"`
	Message     string `json:"message,omitempty"`
	FileName    string `json:"filename,omitempty"`
	S3Name      string `json:"s3name,omitempty"`
	Material    string `json:"material,omitempty"`
	IPAddress   string `json:"IPAddress,omitempty"`
	UploadInfo  string `json:"UploadInfo,omitempty"`
	UploadURL   string `json:"UploadURL,omitempty"`
	DeviceName  string `json:"DeviceName"`
	DeviceUUID  string `json:"DeviceUUID,omitempty"`
	MessageUUID string `json:"MeviceUUID,omitempty"`
	ResultURL   string `json:"ResultURL,omitempty"`
}

// JobClient struct represent the struct that will be sent to devices when sending them a job
type JobClient struct {
	FileName string `json:"filename"`
	Material string `json:"material"`
}

// Config struct represents the configurable values for the Service
type Config struct {
	NumberOfRetries           int
	InitialTimeBetweenRetries int
}

// DLQ_Message struct represent the messages that will be inserted and read from the
// Dead Letter Queue
type DLQ_Message struct {
	Type           string `json:"type"`
	AdditionalInfo string `json:"AdditionalInfo"`
	DeviceName     string `json:"DeviceName"`
	LastResult     string `json:"LastResult"`
	Timestamp      int64  `json:"Timestamp"`
}
