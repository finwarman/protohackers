package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const TCP_PORT = 25565

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", TCP_PORT))
	if err != nil {
		fmt.Println("listen: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("listening on port %d\n", TCP_PORT)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("connection from ", conn.RemoteAddr())
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		fmt.Println("copy: ", err.Error())
	}
}
