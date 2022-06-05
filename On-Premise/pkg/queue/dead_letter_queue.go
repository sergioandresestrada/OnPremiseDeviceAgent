package queue

import "github.com/aws/aws-sdk-go-v2/service/sqs/types"

// DeadLetterQueue interface defines the methods that DeadLetterQueue implementations will need to have
// Iterface is used although only one implementation is used so that we can mock it
type DeadLetterQueue interface {
	SendMessage(string) error
	ReceiveMessages() []types.Message
	RemoveMessage(types.Message) error
}
