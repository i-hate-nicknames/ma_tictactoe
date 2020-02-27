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

// read one client message data from the given reader, parse it
// and return as a message struct
func ReadMessage(reader *bufio.Reader) (interface{}, error) {
	data, err := reader.ReadString('\n')
	if err == io.EOF {
		return nil, fmt.Errorf("client disconnected")
	}
	if err != nil {
		return nil, fmt.Errorf("error reading from client %s", err)
	}
	data = strings.Trim(data, "\n")
	message, err := UnmarshalMessage(data)
	if err != nil {
		return nil, fmt.Errorf("error when parsing client message: %s", err)
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
