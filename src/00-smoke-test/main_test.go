package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

const TEST_TCP_PORT = 44444

var TEST_CONNECT_STR = fmt.Sprintf("localhost:%d", TEST_TCP_PORT)

func TestEchoServer(t *testing.T) {
	go StartServer(TEST_TCP_PORT)      // Start the server in a goroutine
	time.Sleep(time.Millisecond * 500) // (Wait for the server to start)

	// Connect to the server
	conn, err := net.Dial("tcp", TEST_CONNECT_STR)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send a message
	message := "Hello, server!\n"
	fmt.Fprintf(conn, message)

	// Read the response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	// Check if the response matches the message
	if response != message {
		t.Fatalf("Expected '%s', got '%s'", message, response)
	}
}
