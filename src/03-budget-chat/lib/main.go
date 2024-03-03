package budgetchat

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Default tcp port for server
const TCP_PORT = 25565

// Character to indicate sent message is terminated
const MSG_TERM = "\n"

// Entry point
// func main() {
// 	StartServer(TCP_PORT)
// }

// Run the server
func StartServer(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println(S_PREFIX+"listen: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf(S_PREFIX+"listening on port %d\n", port)

	generator := NewIDGenerator()
	broadcaster := NewBroadcaster()

	// Create a goroutine with a connection handler,
	// for each new connection. (Must handle at least 5)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(S_PREFIX+"listener accept error: ", err.Error())
			os.Exit(1)
		}

		fmt.Println(S_PREFIX+"connection from ", conn.RemoteAddr())

		// Create a Client object to represent this connection
		client := &Client{
			id:       int(generator.NextID()),
			joined:   false,
			username: "",
			msgChan:  make(chan Message),
		}

		go HandleConnection(conn, broadcaster, client)
	}
}

func HandleConnection(conn net.Conn, broadcaster *Broadcaster, client *Client) {
	// Close connection and unsubscribe client on disconnect.
	defer func() {
		// Broadcast a leaving message
		if client.joined {
			_, _ = broadcaster.Broadcast(Message{data: fmt.Sprintf("* %s has left the room", client.username)}, client)
		}

		_, _ = broadcaster.Unsubscribe(client)
		conn.Close()
	}()

	// Initial connection message: get username
	welcomeMsg := fmt.Sprintf("[id: %d] Welcome to fubChat! What is your username?", client.id)
	if _, err := conn.Write([]byte(welcomeMsg + MSG_TERM)); err != nil {
		fmt.Println(S_PREFIX+"write error:", err.Error())
		return
	}

	// Buffer for storing received data
	reader := bufio.NewReader(conn)

	// Initial message is username
	for !client.joined {
		// Read data until newline character
		usernameInput, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(S_PREFIX+"read error:", err.Error())
			}
			break
		}

		// Trim newline character
		usernameInput = strings.TrimSuffix(usernameInput, "\n")

		// [Debug] Print received data to STDOUT
		fmt.Printf(S_PREFIX+"received username: '%s'\n", usernameInput)

		if usernameInput != "" {
			if !IsValidUsername(usernameInput) {
				fmt.Printf("%sClient #%d: Invalid username '%s'\n", S_PREFIX, client.id, usernameInput)
				return
			}
			client.username = usernameInput
			client.joined = true
		}
	}

	if client.username == "" {
		fmt.Printf("got empty username for client #%d\n", client.id)
		return
	}

	// Start the message processing goroutine
	go client.ProcessMessages(conn)

	// Subscribe client to receive broadcasts
	_, _ = broadcaster.Subscribe(client)

	// Broadcast '* The room contains: ...' to newly joined user
	var otherClients []string
	for _, c := range broadcaster.clients {
		if client.id != c.id && c.username != "" {
			otherClients = append(otherClients, c.username)
		}
	}
	otherClientsStr := strings.Join(otherClients, ", ")
	usersMessage := "* The room contains: " + otherClientsStr
	client.QueueMessage(Message{data: usersMessage})

	// Broadcast 'joined' message to all connected clients
	msg := Message{data: fmt.Sprintf("* %s has entered the room", client.username)}
	_, _ = broadcaster.Broadcast(msg, client)

	// While connection is open, check for data to read
	for {
		// Read data until newline character
		data, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(S_PREFIX+"read error:", err.Error())
			}
			break
		}

		// Trim newline character
		data = strings.TrimSuffix(data, "\n")

		// [Debug] Print received data to STDOUT
		fmt.Printf(S_PREFIX+"received: '%s'\n", data)

		messageStr := fmt.Sprintf("[%s] %s", client.username, data)
		msg := Message{
			data: messageStr,
		}

		_, _ = broadcaster.Broadcast(msg, client)
	}
}

// Example client message processing loop
func (c *Client) ProcessMessages(conn net.Conn) {
	// Continually read from message channel
	for msg := range c.msgChan {
		if len(msg.data) > 0 {
			fmt.Printf("%s(@%s) sending message: '%s'\n", S_PREFIX, c.username, msg.data)

			// Send the response
			if _, err := conn.Write([]byte(msg.data + MSG_TERM)); err != nil {
				fmt.Println(S_PREFIX+"write error:", err.Error())
				break
			}
		}
	}
}
