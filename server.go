package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func startServer(port string, done chan<- bool) {
	serverSock, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start a server: %s\n", err)
	}
	done <- true
	for {
		conn, err := serverSock.Accept()
		if err != nil {
			log.Printf("Failed to handle a client: %s\n", err)
			continue
		}
		go handleClient(conn)
	}
}

// todo: implement server struct that will hold connection instead of passing it around
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
		data = strings.Trim(data, "\n")
		message, err := UnmarshalMessage(data)
		if err != nil {
			log.Printf("Error when parsing client message: %s", err)
		} else {
			handleMessage(conn, message)
		}
	}
}

func handleMessage(conn net.Conn, message interface{}) {
	var data string
	var err error
	switch message := message.(type) {
	case QuestionMessage:
		answer := AnswerMessage{"I have no idea about " + message.Text, "have a nice life"}
		data, err = MarshalMessage(answer)
	case AnswerMessage:
		answer := AnswerMessage{"nothing", "you are not allowed to answer questions son"}
		data, err = MarshalMessage(answer)
	default:
		log.Printf("Unsupported message type: %T", message)
	}
	if err != nil {
		log.Printf("Error when marshaling message: %s", err)
	}
	fmt.Fprintln(conn, data)
}
