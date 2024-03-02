package main

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
	// TODO: actually handle message stuff
	// fmt.Printf("%s(@%s) queueing message: '%s'\n", S_PREFIX, c.username, message.data)
	fmt.Printf("(@%s) queueing message: '%s'\n", c.username, message.data)

	// Append the message to this client's message channel
	c.msgChan <- message
}
