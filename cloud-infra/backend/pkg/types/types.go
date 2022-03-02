package types

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
