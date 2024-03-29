package main

import (
	"fmt"
	"io"
	"net"
	"os"
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

	if _, err := io.Copy(conn, conn); err != nil {
		fmt.Println("copy: ", err.Error())
	}
}
