package queue

import "github.com/aws/aws-sdk-go-v2/service/sqs/types"

type Queue interface {
	ReceiveMessages() []types.Message
	RemoveMessage(types.Message) error
}
