package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// connect to the server, send string, receive a single response
// and print it
func startClient(addr string) {
	log.Println("Starting client, connecting to " + addr)
	conn, err := net.Dial("tcp4", addr)
	defer conn.Close()
	if err != nil {
		log.Fatalf("Couldn't connect to %s", addr)
		return
	}
	sockReader := bufio.NewReader(conn)
	for {
		message := readServerMessage(sockReader)
		handleServerMessage(message)
	}
}

func readInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')
	if err == io.EOF {
		fmt.Println("Exiting")
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Error reading userInput: %s\n", err)
		os.Exit(1)
	}
	return strings.Trim(input, "\n")
}

// todo: move to client package and rename to readMessage
func readServerMessage(reader *bufio.Reader) interface{} {
	serverData, err := reader.ReadString('\n')
	if err == io.EOF {
		fmt.Println("Server closed the connection")
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Error reading reply from server: %s\n", err)
		os.Exit(1)
	}
	serverData = strings.Trim(serverData, "\n")
	message, err := UnmarshalMessage(serverData)
	if err != nil {
		fmt.Printf("Error unmarshaling server message: %s\n", err)
		os.Exit(1)
	}
	return message
}

// todo: move to client package and rename to handleMessage
func handleServerMessage(message interface{}) {
	switch message := message.(type) {
	case WaitingMessage:
		fmt.Println("Waiting for another player to connect")
	case BoardMessage:
		fmt.Printf("Server sent us a board!\n%s\n", message.board)
	default:
		log.Printf("Unsupported message type: %T", message)
	}
}
