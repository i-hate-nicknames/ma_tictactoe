package messaging

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// todo check if this is a blocking call
func SendMessage(conn net.Conn, message interface{}) {
	data, err := MarshalMessage(message)
	if err != nil {
		log.Printf("Error marshaling message: %s\n", err)
		return
	}
	fmt.Fprintln(conn, data)
}

func ReadMessages(conn net.Conn, messages chan<- interface{}, errors chan<- error) {
	reader := bufio.NewReader(conn)
	for {
		message, err := ReadMessage(reader)
		if err != nil {
			errors <- err
			return
		}
		messages <- message
	}
}

func ReadMessage(reader *bufio.Reader) (interface{}, error) {
	data, err := reader.ReadString('\n')
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	data = strings.Trim(data, "\n")
	message, err := UnmarshalMessage(data)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// todo: refactor both read message method into one, possibly by introducing
// channels
func ReadServerMessage(reader *bufio.Reader) interface{} {
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
