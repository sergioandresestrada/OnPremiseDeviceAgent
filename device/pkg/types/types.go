package types

// JobDevice defines the struct that the JSON information received
// by the POST /job endpoint should receive.
type JobDevice struct {
	FileName string `json:"filename"`
	Material string `json:"material"`
}
