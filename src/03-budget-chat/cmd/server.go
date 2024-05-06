package main

import budgetchat "github.com/finwarman/protohackers/budgetchat/lib"

const TCP_PORT = budgetchat.DEFAULT_TCP_PORT

func main() {
	// Start the server
	budgetchat.StartServer(TCP_PORT)
}
