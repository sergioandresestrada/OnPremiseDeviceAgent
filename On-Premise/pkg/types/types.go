package types

type Message struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	FileName  string `json:"filename,omitempty"`
	S3Name    string `json:"s3name,omitempty"`
	Material  string `json:"material,omitempty"`
	IPAddress string `json:"IPAddress,omitempty"`
}

type JobClient struct {
	FileName string `json:"filename"`
	Material string `json:"material"`
}
