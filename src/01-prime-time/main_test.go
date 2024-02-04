package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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

	// List of tests of format message: {"method":"isPrime","number":19809878}
	// response: "true", "false", "malformed"
	tests := []struct {
		message  string
		expected string
	}{
		// Data to Send, Expected Response Type
		{"{\"method\":\"isPrime\",\"number\":13441}\n", "true"},
		{"{\"method\":\"isPrime\",\"number\":12659}\n", "true"},
		{"{\"method\":\"isPrime\",\"number\":12241}\n", "true"},

		{"{\"method\":\"isPrime\",\"number\":123456}\n", "false"},
		{"{\"method\":\"isPrime\",\"number\":-1}\n", "false"},
		{"{\"method\":\"isPrime\",\"number\":-13441}\n", "false"},
		{"{\"method\":\"isPrime\",\"number\":7.123}\n", "false"},
		{"{\"method\":\"isPrime\",\"number\":13441.123}\n", "false"},

		{"{\"method\":\"isPrime\",\"number\":\"abc\"}\n", "malformed"},
		{"{\"method\":\"something\",\"number\":\"abc\"}\n", "malformed"},
		{"{\"method\":\"isPrime\",\"number\":\"123\"}\n", "malformed"},
		{"{method, isPrime, number, 13441}\n", "malformed"},
	}

	for _, test := range tests {
		// Send a message
		fmt.Fprint(conn, test.message)

		// Read the response
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read from connection: %v", err)
		}

		// Check if the response matches the expected value
		if test.expected == "malformed" {
			if !strings.Contains(response, "[]") {
				t.Fatalf("Expected malformed message, got '%s'", response)
			}
		} else if !strings.Contains(response, test.expected) {
			fmt.Printf("Sent data: %s\n", test.message)
			t.Fatalf("Expected to contain '%s', got '%s'", test.expected, response)
		}
	}
}
