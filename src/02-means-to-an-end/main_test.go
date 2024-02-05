package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
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

	insertExampleBytes, _ := hexStringToByteArray("" +
		// Hexadecimal                     Decoded
		"49 00 00 30 39 00 00 00 65 " + // I 12345 101
		"49 00 00 30 3a 00 00 00 66 " + // I 12346 102
		"49 00 00 30 3b 00 00 00 64 " + // I 12347 100
		"49 00 00 a0 00 00 00 00 05 " + // I 40960 5

		// Query
		"51 00 00 30 00 00 00 40 00 ", // Q 12288 16384
	)

	// Expected response:
	// 00 00 00 65
	// (101)

	reader := bufio.NewReader(conn)

	if _, err := conn.Write(insertExampleBytes); err != nil {
		t.Fatalf("Failed to send bytes: %v", err)
	}

	// Read the response
	responseBytes := make([]byte, 4)
	count, err := reader.Read(responseBytes)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	response := string(responseBytes)

	fmt.Printf("[test client] received response (%d bytes): %s\n",
		count, strconv.Quote(string(response)))

	expected, _ := hexStringToByteArray("00 00 00 65")
	if response != string(expected) {
		t.Fatalf(
			"Expected response does not match, expected '%X', got '%X'",
			expected, response,
		)
	}

	// Time to allow output to flush
	time.Sleep(time.Millisecond * 100)
}

func hexStringToByteArray(hexStr string) ([]byte, error) {
	splits := strings.Fields(hexStr)
	byteArray := make([]byte, len(splits))

	// Convert each split into a byte and add to the byte array
	for i, s := range splits {
		b, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			return nil, err
		}
		byteArray[i] = byte(b)
	}

	return byteArray, nil
}
