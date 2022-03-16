package queue

// Queue interface defines the methods that Queue implementations will need to have
// Iterface is used although only one implementation is used so that we can mock it
type Queue interface {
	SendMessage(string) error
}
