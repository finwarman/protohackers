package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

// =================================
// == Problem 4: Unusual Database ==
//
// Key-Value store accessed over UDP
//
// Each request and each response is a single UDP packet, of up to 1000 bytes.
//
// two types of request: INSERT and RETRIEVE
//  - INSERT:   insert a value for a key      (Contains an equals sign, '=')
//  - RETRIEVE: retrevie the value of a key   (Does NOT contain an equals sign)
//
// INSERT does not yield a response. RETRIEVE returns '{key}={value}'.
//
// Key 'version' returns the version string, this key cannot be modified.
//
// =================================

//
// === CONSTANTS === //
//

// Prefix for colourised server log messages
const (
	ColourReset = "\033[0m"
	ColourRed   = "\033[31m"
	ColourCyan  = "\033[36m"
)
const S_PREFIX = ColourCyan + "[server]" + ColourReset + " "
const S_ERROR = ColourRed + "[server]" + ColourReset + " "

const UDP_PORT = 25565

const MAX_REQ_BYTES = 1000
const MAX_RES_BYTES = 1000

const VERSION = "FunkyDatabase@v1.0.0"

// key-balue store with a mutex for safe concurrent access
var DATABASE = make(map[string]string)
var DB_MUTEX sync.RWMutex

//
// === METHODS === //
//

func main() {
	StartServer(UDP_PORT)
}

func StartServer(port int) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)

	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Printf(S_PREFIX+"listen error: on UDP port %d: %v\n", port, err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("%sUDP: listening on port %d\n", S_PREFIX, port)

	buffer := make([]byte, MAX_REQ_BYTES)

	DATABASE["version"] = VERSION

	// Handle incoming packets
	for {
		n, remoteAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println(S_PREFIX+"read error: ", err.Error())
			continue
		}

		// Process the received packet
		handlePacket(conn, remoteAddr, buffer[:n])
	}
}

func handlePacket(conn net.PacketConn, addr net.Addr, data []byte) {
	input := string(data) // NOTE: don't remove newlines from datagram!
	fmt.Printf(S_PREFIX+"received from %s: %s\n",
		addr.String(), strconv.Quote(input))

	switch {
	case strings.Contains(input, "="):
		// INSERT
		parts := strings.SplitN(input, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			// Prevent modification of the 'version' key
			if key == "version" {
				fmt.Println(S_PREFIX + "[INSERT] Update to 'version' DENIED")
				return
			}

			fmt.Printf(S_PREFIX+"[INSERT]: Set key %s to value %s\n",
				strconv.Quote(key), strconv.Quote(value))

			DB_MUTEX.Lock()
			defer DB_MUTEX.Unlock()
			DATABASE[key] = value
		}
	default:
		// RETRIEVE
		// This also handles empty datagrams, returning '='
		key := input
		value, exists := DATABASE[key]
		if !exists {
			value = ""
		}

		fmt.Printf(S_PREFIX+"[RETRIEVE]: key %s has value %s\n",
			strconv.Quote(key), strconv.Quote(value))

		res := fmt.Sprintf("%s=%s", key, value)
		sendResponse(conn, addr, res)
	}
}

func sendResponse(conn net.PacketConn, addr net.Addr, res string) {
	data := []byte(res)

	// Truncate reposnse before sending
	if len(data) > MAX_RES_BYTES {
		data = data[:MAX_RES_BYTES]
	}

	_, err := conn.WriteTo(data, addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}
