package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

const TEST_TCP_PORT = 44444

// prefix for client log messages
const C_PREFIX = ColorYellow + "[client]" + ColorReset + " "

var TEST_CONNECT_STR = fmt.Sprintf("localhost:%d", TEST_TCP_PORT)

func TestEchoServer(t *testing.T) {
	go StartServer(TEST_TCP_PORT)      // Start the server in a goroutine
	time.Sleep(time.Millisecond * 500) // (Wait for the server to start)

	go StartNewClient(t, 1)

	time.Sleep(time.Millisecond * 50)

	go StartNewClient(t, 2)

	// Time to allow output to flush
	time.Sleep(time.Millisecond * 500)
}

func StartNewClient(t *testing.T, id int) {
	// Connect to the server
	conn, err := net.Dial("tcp", TEST_CONNECT_STR)
	if err != nil {
		t.Fatalf(C_PREFIX+"failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Buffer for storing received data
	reader := bufio.NewReader(conn)

	// Await username message
	msg, err := reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			fmt.Println(C_PREFIX+"read error:", err.Error())
		}
	}
	msg = strings.TrimSuffix(msg, "\n")

	fmt.Printf(C_PREFIX+"received: '%s'\n", msg)

	// Delay
	time.Sleep(time.Millisecond * 100)

	// Send username
	username := "username" + fmt.Sprintf("%d", id)
	if _, err := conn.Write([]byte(username + "\n")); err != nil {
		t.Fatalf(C_PREFIX+"failed to send bytes: %v", err)
	}

	// Delay
	time.Sleep(time.Millisecond * 500)

	// Send a message
	message := "this is my message " + fmt.Sprintf("%d", id)
	if _, err := conn.Write([]byte(message + "\n")); err != nil {
		t.Fatalf(C_PREFIX+"failed to send bytes: %v", err)
	}
}

// TODO - simulate client message processing loop (recieve/send)
