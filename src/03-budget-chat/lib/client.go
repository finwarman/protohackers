package budgetchat

import "fmt"

// Message represents a client message
type Message struct {
	data string
}

// Client represents a chat client
type Client struct {
	id       int
	joined   bool
	username string
	msgChan  chan Message // message queue
}

// QueueMessage adds a message to the message queue for a client
func (c *Client) QueueMessage(message Message) {
	fmt.Printf(S_PREFIX+"(@%s) queueing message: '%s'\n", c.username, message.data)

	// Push message to this client's message channel
	c.msgChan <- message
}
