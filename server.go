package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func startServer(port string) {
	serverSock, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start a server: %s\n", err)
	}
	go func() {
		for {
			conn, err := serverSock.Accept()
			if err != nil {
				log.Printf("Failed to handle a client: %s\n", err)
				continue
			}
			go handleClient(conn)
		}
	}()
}

// handle a client: reply to every message with modified client message
func handleClient(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Printf("Client disconnected")
			return
		}
		if err != nil {
			log.Printf("Error reading from client %s\n", err)
			return
		}
		fmt.Fprintf(conn, data+" yourself")
	}
}
