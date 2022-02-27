package queue

type Queue interface {
	SendMessage(string) error
}
