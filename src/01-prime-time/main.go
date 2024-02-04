package main

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"strings"

	"github.com/finwarman/protohackers/src/lib/json"
)

const TCP_PORT = 25565

func main() {
	StartServer(TCP_PORT)
}

func StartServer(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("listen: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("listening on port %d\n", port)

	// Create a goroutine with a connection handler,
	// for each new connection. (Must handle at least 5)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("connection from ", conn.RemoteAddr())

		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// Buffer for storing received data
	reader := bufio.NewReader(conn)

	// While connection is open, check for data to read
	for {
		// Read data until newline character
		data, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err.Error())
			}
			break
		}

		// Trim newline character
		data = strings.TrimSuffix(data, "\n")

		// [Debug] Print received data to STDOUT
		fmt.Printf("received: %s\n", data)

		// Handle JSON request
		response := handleJSON(data)
		fmt.Printf("sending response: %s\n", response)

		// Send response, terminated with newline
		if _, err := conn.Write([]byte(string(response) + "\n")); err != nil {
			fmt.Println("write error:", err.Error())
			break
		}
	}
}

var MALFORMED_RESPONSE []byte = []byte("[]")

// Validate JSON request:
// Input must:
//   - Be valid JSON
//   - Have `/method` = "isPrime"
//   - Type of `/number` is number
//   - Extraneous fields are ignored
//
// Uses the JSON parser `github.comfinwarman/protohacker/src/lib/json` -
// a minimal parser written as a learning exercise for this project.
//
// Response format:
//
//	{"method":"isPrime","prime":false}
//
// If request is malformed, send a malformed response
//
//	e.g. '[]'
func handleJSON(data string) []byte {
	parsedValue, err := json.ParseJSON(data)
	if err != nil {
		fmt.Printf("parsing error:\n%v\n\ninput data:%s\n\n", err, data)
		return MALFORMED_RESPONSE
	}

	fmt.Printf("parsed JSON value:\n%v\n", parsedValue)

	parsedValueMapGeneric := json.ConvertToNative(parsedValue)

	// Convert to `string: object`
	parsedValueMap, ok := parsedValueMapGeneric.(map[string]interface{})
	if !ok {
		fmt.Println("incorrect type or invalid object ")
		return MALFORMED_RESPONSE
	}

	// get parsedValueMap["method"] string
	method, ok := parsedValueMap["method"].(string)
	if !ok {
		fmt.Println("`/method` not found or not type string")
		return MALFORMED_RESPONSE
	} else {
		fmt.Printf("got method: %s\n", method)
	}

	// validate method used
	if method != "isPrime" {
		fmt.Println("`/method` was not `isPrime`")
		return MALFORMED_RESPONSE
	}

	// wrong number format, but not malformed
	float, ok := parsedValueMap["number"].(float64)
	if ok {
		fmt.Printf("`/number` was a float %f, not int\n", float)
		return []byte("{\"method\":\"isPrime\",\"prime\":false}")
	}

	// get parsedValueMap["number"] int
	number, ok := parsedValueMap["number"].(int)
	if !ok {
		fmt.Println("`/number` not found or not type number (int)")
		return MALFORMED_RESPONSE
	} else {
		fmt.Printf("got number: %d\n", number)
	}

	// return properly-formed response
	isPrime := big.NewInt(int64(number)).ProbablyPrime(0)
	returnJSONStr := fmt.Sprintf("{\"method\":\"isPrime\",\"prime\":%t}", isPrime)

	return []byte(returnJSONStr)
}
