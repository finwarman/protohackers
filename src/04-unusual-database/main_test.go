package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

const TEST_UDP_PORT = 44444

var TEST_CONNECT_STR = fmt.Sprintf("localhost:%d", TEST_UDP_PORT)

func TestEchoServer(t *testing.T) {
	go StartServer(TEST_UDP_PORT)      // Start the server in a goroutine
	time.Sleep(time.Millisecond * 500) // Wait for the server to start

	conn, err := net.Dial("udp", TEST_CONNECT_STR)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// === INITIAL INSERT AND RETRIEVE ===

	// Send an 'INSERT' to the server
	message := "foo=bar=baz"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Retrieve the inserted value
	message = "foo"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Receive a response
	responseBytes := make([]byte, 1000)
	n, err := conn.Read(responseBytes)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}
	response := string(responseBytes[:n])

	// Check the echo response
	message = "foo=bar=baz"
	if response != message {
		t.Fatalf("Expected response '%s', got '%s'", message, response)
	}

	// === UPDATE AND RETRIEVE ==

	// Send an update 'INSERT' to the server
	message = "foo="
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Retrieve the updated value
	message = "foo"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	responseBytes = make([]byte, 1000)
	n, err = conn.Read(responseBytes)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}
	response = string(responseBytes[:n])

	// Check the echo response
	message = "foo="
	if response != message {
		t.Fatalf("Expected response '%s', got '%s'", message, response)
	}

	// === VERSION: UPDATE AND REQUEST ===

	// Attempt to update 'VERSION' (ignored)
	message = "version=FOOBAR"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Retrieve the version value
	message = "version"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	responseBytes = make([]byte, 1000)
	n, err = conn.Read(responseBytes)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}
	response = string(responseBytes[:n])

	message = "version=FunkyDatabase@v1.0.0"
	if response != message {
		t.Fatalf("Expected response '%s', got '%s'", message, response)
	}

	// === EDGE CASES === //

	// Test updating strings with trailing newlines and whitespace
	// Send an 'INSERT' with trailing newline and whitespace
	message = "testKey=testValue\n \t\r\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Retrieve the inserted value
	message = "testKey"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Receive a response
	responseBytes = make([]byte, 1000)
	n, err = conn.Read(responseBytes)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}
	response = string(responseBytes[:n])

	// Check the echo response
	message = "testKey=testValue\n \t\r\n"
	if response != message {
		t.Fatalf("Expected response '%s', got '%s'", message, response)
	}

	// Test empty datagram returns '='
	_, err = conn.Write([]byte(""))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Receive a response for the empty datagram
	responseBytes = make([]byte, 1000)
	n, err = conn.Read(responseBytes)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}
	response = string(responseBytes[:n])

	// Check the echo response for empty datagram
	message = "="
	if response != message {
		t.Fatalf("Expected response '%s', got '%s'", message, response)
	}

	// ===================================

	// Time to allow output to flush
	time.Sleep(time.Millisecond * 500)
}
