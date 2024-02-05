package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"strconv"
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

	// Database for this client
	// (Each client has a different asset)
	assetDatabase := make(map[int32]int32)

	// While connection is open, check for data to read
	for {
		// Buffer to hold exactly 9 bytes
		const messageSize = 9
		buf := make([]byte, messageSize)

		// Read exactly 9 bytes
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection closed by client")
			} else {
				fmt.Println("read error:", err.Error())
			}
			break
		}

		// [Debug] Print received data to STDOUT
		// fmt.Print("received (hex): ")
		// for _, b := range buf[:] {
		// 	fmt.Printf("%02x ", b)
		// }
		// fmt.Println()

		// Handle the 9-byte message
		response := handleBytesData(buf, &assetDatabase)

		if len(response) > 0 {
			fmt.Printf("sending response: %s\n",
				strconv.Quote(string(response)))

			// Send the response
			if _, err := conn.Write([]byte(response)); err != nil {
				fmt.Println("write error:", err.Error())
				break
			}
		}
	}
}

/*
Each request is exactly 9 bytes
It is either an INSERT or a QUERY

Message format:
Byte:  |  0  |  1     2     3     4  |  5     6     7     8  |
Type:  |char |         int32         |         int32         |

Char indicates message `type`, either `I` or `Q` (INSERT or QUERY).

The next 8 bytes are two signed two's complement 32-bit integers
in network byte order (big endian),

For INSERT:
  - First int32 is timestamp, in seconds since 00:00, 1st Jan 1970.
  - Second int32 is price, in pennies, of this client's asset, at the
    given timestamp.

For QUERY:
  - First int32 is mintime, the earliest timestamp of the period.
  - Second int32 is maxtime, the latest timestamp of the period.

The server must compute the mean of the inserted prices with timestamps T,
mintime <= T <= maxtime (i.e. timestamps in closed interval [mintime, maxtime]).

If the mean is not an integer, it is acceptable to round either up or down,
at the server's discretion.

The server must then send the mean to the client as a single int32.

If there are no samples within the requested period, or if mintime comes
after maxtime, the value returned must be 0.

Behaviour is undefined if there are multiple prices with the same timestamp from
the same client.

Where a client triggers undefined behaviour, the server can do anything it likes
for that client, but must not adversely affect other clients that did not
trigger undefined behaviour.
*/
func handleBytesData(data []byte, assetDatabase *map[int32]int32) []byte {

	var UNDEF_RESPONSE = []byte("undef\n")

	// Parse and validate operation-type byte char
	charByte := data[0]
	if charByte != 'I' && charByte != 'Q' {
		fmt.Printf("invalid char byte: %02x\n", charByte)
		return UNDEF_RESPONSE
	}

	// Convert bytes 1-4 to a signed 32-bit integer
	intOneValue := int32(binary.BigEndian.Uint32(data[1:5]))
	// fmt.Printf("1st int32 value is: %d\n", intOneValue)

	// Convert bytes 1-4 to a signed 32-bit integer
	intTwoValue := int32(binary.BigEndian.Uint32(data[5:9]))
	// fmt.Printf("2nd int32 value is: %d\n", intTwoValue)

	// Handle INSERT
	// - Insert a timestamped price
	if charByte == 'I' {
		// fmt.Println("command is I for INSERT")

		timestamp := intOneValue
		price := intTwoValue

		// fmt.Printf("inserting price $%d at time %d\n", price, timestamp)

		(*assetDatabase)[timestamp] = int32(price)

		return nil
	}

	// Handle QUERY
	// - Fetch a mean price across period
	if charByte == 'Q' {
		// fmt.Println("command is Q for QUERY")

		mintime := intOneValue
		maxtime := intTwoValue

		// fmt.Printf("querying time range %d - %d\n", mintime, maxtime)

		total := int64(0)
		count := 0
		for timestamp, price := range *assetDatabase {
			if timestamp >= mintime && timestamp <= maxtime {
				count++
				total += int64(price)
			}
		}

		mean := int64(0)
		if count > 0 {
			mean = int64(math.Round(float64(total) / float64(count)))
		}

		// fmt.Printf("got mean in timerange: %d (%d/%d)\n", mean, total, count)

		meanBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(meanBytes, uint32(int32(mean)))

		return meanBytes
	}

	return UNDEF_RESPONSE
}
