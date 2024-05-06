package budgetchat

import (
	"fmt"
	"net"
)

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
	// usernameMsg := Colourise("@"+c.username, ColourYellow)
	// fmt.Printf("%s%s\tqueueing message: '%s'\n", S_PREFIX, usernameMsg, message.data)

	// Push message to this client's message channel
	c.msgChan <- message
}

// ProcessMessages is the client message processing loop
func (c *Client) ProcessMessages(conn net.Conn) {
	// Continually read from message channel
	for msg := range c.msgChan {
		if len(msg.data) > 0 {
			usernameMsg := Colourise("@"+c.username, ColourYellow)
			fmt.Printf("%s%s\tsending message: '%s'\n", S_PREFIX, usernameMsg, msg.data)

			// Send the response
			if _, err := conn.Write([]byte(msg.data + MSG_TERM)); err != nil {
				fmt.Println(S_PREFIX+"write error:", err.Error())
				break
			}
		}
	}
}
