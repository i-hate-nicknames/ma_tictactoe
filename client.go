package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// connect to the server, send string, receive a single response
// and print it
func startClient(addr string) {
	log.Println("Starting client, connecting to " + addr)
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		log.Fatalf("Couldn't connect to %s", addr)
		return
	}
	fmt.Fprintf(conn, "kurwa\n")
	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading server reply %s\n", err)
	} else {
		fmt.Println("Server replied " + data)
	}
}
