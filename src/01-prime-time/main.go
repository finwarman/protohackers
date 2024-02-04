package main

import (
	"fmt"
	"io"
	"math/big"
	"net"
	"os"

	"github.com/finwarman/protohackers/src/lib/json"
)

const TCP_PORT = 25565

func main() {
	StartServer(TCP_PORT)
}

func StartServer(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
	buf := make([]byte, 32768)

	// While connection is open, check for data to read
	for {
		// If there is data to read, place it into the buffer
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err.Error())
			}
			break
		}

		// [Debug] Print recieved data to STDOUT
		data := buf[:n]
		fmt.Printf("received: %s\n", string(data))

		response := handleJSON(string(data))
		fmt.Printf("sending response: %s\n", response)

		// Echo the data back to the client
		if _, err := conn.Write(response); err != nil {
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
// Uses the JSON parser `github.comfinwarman/prothacker/`
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
	returnJSONStr := fmt.Sprintf("{\"method\":\"isPrime\",\"prime\":%t}\n", isPrime)

	return []byte(returnJSONStr)
}
