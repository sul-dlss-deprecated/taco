package streaming

// Message represents the message we send to kinesis
type Message struct {
	Action string
	ID     string
}
