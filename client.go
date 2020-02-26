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
		handleServerMessage(message, conn)
	}
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
func handleServerMessage(message interface{}, conn net.Conn) {
	switch message := message.(type) {
	case WaitingMessage:
		fmt.Println("Waiting for another player to connect")
	case BoardMessage:
		board := message.Board
		// print state
		fmt.Printf("Server sent us a board!\n%s\n", board)
		reply, err := readMove(board)
		if err != nil {
			fmt.Println(err)
			// retry handling
			handleServerMessage(message, conn)
		}
		sendClientMessage(conn, reply)
	case ErrorMessage:
		fmt.Printf("Error: %s\n", message.Text)
	case HelloMessage:
		fmt.Printf("%s\nYour player is %s\n", message.Text, message.AssignedPlayer)
	default:
		log.Printf("Unsupported message type: %T", message)
	}
}

func readMove(board *Board) (MoveMessage, error) {
	var result MoveMessage
	var x, y int
	fmt.Println("Enter x coordinate")
	_, err := fmt.Scanf("%d\n", &x)
	if err != nil {
		return result, err
	}
	fmt.Println("Enter y coordinate")
	_, err = fmt.Scanf("%d\n", &y)
	if err != nil {
		return result, err
	}
	err = board.validateCoordinates(x, y)
	if err != nil {
		return result, err
	}
	return MoveMessage{x, y}, nil
}

func readInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err == io.EOF {
		fmt.Println("Exiting")
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Error reading userInput: %s\n", err)
		return "", err
	}
	return strings.Trim(input, "\n"), nil
}

func sendClientMessage(conn net.Conn, message interface{}) {
	data, err := MarshalMessage(message)
	if err != nil {
		log.Printf("Error marshaling message: %s\n", err)
		return
	}
	fmt.Fprintln(conn, data)
}
