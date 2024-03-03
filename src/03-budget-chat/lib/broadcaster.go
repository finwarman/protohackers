package budgetchat

import "fmt"

// Broadcaster handles sending messages to multiple clients
type Broadcaster struct {
	clients map[int]*Client // id -> client
}

// NewBroadcaster creates a new Broadcaster instance
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[int]*Client),
	}
}

// Subscribe adds a client to the broadcaster
func (b *Broadcaster) Subscribe(client *Client) (bool, error) {
	if client == nil || client.id == 0 {
		return false, fmt.Errorf("subscription failed: no client provided")
	}

	// Subscribe client (overwrites if already subscribed)
	b.clients[client.id] = client

	return true, nil
}

// Unsubscribe removes a client from the broadcaster
func (b *Broadcaster) Unsubscribe(client *Client) (bool, error) {
	if client == nil || client.id == 0 {
		return false, fmt.Errorf("unsubscription failed: no client provided")
	}

	// Unsubscribe client, safely ignore if client was already unsubscribed
	delete(b.clients, client.id)

	return true, nil
}

// Broadcast sends a message to all clients except the source client
func (b *Broadcaster) Broadcast(message Message, sourceClient *Client) (bool, error) {
	skipID := 0
	if sourceClient != nil && sourceClient.id != 0 {
		skipID = sourceClient.id
	}

	for clientID, client := range b.clients {
		// Don't send message back to source client (if set)
		if skipID > 0 && clientID == skipID {
			continue
		}

		client.QueueMessage(message)
	}

	return true, nil
}
