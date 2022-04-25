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

// Information struct represents the names of the available files with information about the devices
// that are present in the object storage
type Information struct {
	Jobs           []string
	Identification []string
}

// Device struct represents the information about a device that we have, readed from the Database or received
// from an API call to create and store a new one
type Device struct {
	DeviceUUID string `json:"DeviceUUID,omitempty"`
	IP         string `json:"IP"`
	Name       string `json:"Name"`
	Model      string `json:"Model,omitempty"`
}

// Response struct represents the information received from the On-Premise server about the outcome of a message
type Response struct {
	Result    string `json:"Result"`
	Timestamp int64  `json:"Timestamp"`
}

// MessageDB struct represents the information about a message that is inserted into the DB
type MessageDB struct {
	DeviceUUID     string
	MessageUUID    string
	Type           string
	AdditionalInfo string
	Timestamp      int64
	// this field is only used to read info from DynamoDB and not sent in JSON responses
	Information string `json:"-"`
}

// MessageDB struct represents the information about a message that is inserted into the DB
type ResultDB struct {
	DeviceUUID  string
	MessageUUID string
	Result      string
	Timestamp   int64
}
