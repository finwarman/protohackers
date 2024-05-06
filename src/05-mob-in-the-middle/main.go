package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const TCP_PORT = 25565
const UPSTREAM_HOST = "chat.protohackers.com"
const UPSTREAM_PORT = 16963

func main() {
	upstream := fmt.Sprintf("%s:%d", UPSTREAM_HOST, UPSTREAM_PORT)
	StartServer(TCP_PORT, upstream)
}

func StartServer(port int, upstream string) {
	fmt.Printf("will forward to upstream: %s\n", upstream)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("listen: ", err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("listening on port %d\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("connection from ", conn.RemoteAddr())

		go HandleConnection(conn, upstream)
	}
}

func HandleConnection(clientConn net.Conn, upstream string) {
	defer clientConn.Close()

	// Connect to upstream server (connection per mitm'd client)
	upstreamConn, err := net.Dial("tcp", upstream)
	if err != nil {
		fmt.Println("dial upstream: ", err.Error())
		return
	}
	defer upstreamConn.Close()

	// Channel to communicate closure of either connection
	// (Disconnect both ends on a closed client)
	closed := make(chan bool, 1)

	// Async handlers for client and upstream, with message rewriting:

	// Forward messages from client to upstream
	go func() {
		reader := bufio.NewReader(clientConn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("read from client: ", err.Error())
				break
			}
			message = message[:len(message)-1]

			fmt.Println("Received from client:", strconv.Quote(message)) // Print actual message content
			rewrittenMessage := RewriteCoins(message)                    // Rewrite the message
			_, err = io.WriteString(upstreamConn, rewrittenMessage+"\n")
			if err != nil {
				fmt.Println("write to upstream: ", err.Error())
				break
			}
		}
		closed <- true // Signal that client connection closed
	}()

	// Forward messages from upstream to client
	go func() {
		reader := bufio.NewReader(upstreamConn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("read from upstream: ", err.Error())
				break
			}
			message = message[:len(message)-1]

			fmt.Println("Received from upstream:", strconv.Quote(message)) // Print actual message content
			rewrittenMessage := RewriteCoins(message)                      // Rewrite the message
			_, err = io.WriteString(clientConn, rewrittenMessage+"\n")
			if err != nil {
				fmt.Println("write to client: ", err.Error())
				break
			}
		}
		closed <- true // Signal that upstream connection closed
	}()

	// Wait for either connection to close
	<-closed
}

// Coin Constants
const TARGET_BOGUS_ADDRESS = "7YWHMfk9JZe0LM0g1ZauHuiSxhI" // Tony's Boguscoin address
const DUMMY = "_YWHMfk9JZe0LM0g1ZauHuiSxhI"                // starts with '_' not '7'

var BOGUSCOIN_REGEX = regexp.MustCompile(`(^|[ ])7[0-9a-zA-Z]{25,34}([ ]|$)`)
var REPLACEMENT_REGEX = regexp.MustCompile(`(^ ` + DUMMY + `)|(` + DUMMY + ` $)`)

func RewriteCoins(message string) string {
	if BOGUSCOIN_REGEX.Match([]byte(message)) {
		fmt.Println("[REWRITE (matched coin regex)")
		fmt.Println(" - Original message:", strconv.Quote(message))
		// Hack to ensure consecutive replacements are achieved
		replaced := message
		for i := 0; i < len(message)/26; i++ {
			replaced = BOGUSCOIN_REGEX.ReplaceAllString(replaced, " "+DUMMY+" ")
		}
		replaced = REPLACEMENT_REGEX.ReplaceAllString(replaced, DUMMY)
		replaced = strings.ReplaceAll(replaced, DUMMY, TARGET_BOGUS_ADDRESS)
		fmt.Println(" - Replaced message:", strconv.Quote(string(replaced)))
		return string(replaced)
	}
	return message
}
