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
	inputReader := bufio.NewReader(os.Stdin)
	for {
		userInput, err := readInput(inputReader)
		userInput = strings.Trim(userInput, "\n")
		if err == io.EOF {
			fmt.Println("Exiting")
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("Error reading userInput: %s\n", err)
			os.Exit(1)
		}
		out, err := MarshalMessage(QuestionMessage{userInput})
		if err != nil {
			fmt.Printf("Error marshaling user input: %s\n", err)
			continue
		}
		fmt.Fprintln(conn, out)
		reply, err := sockReader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Server closed the connection")
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("Error reading reply from server: %s\n", err)
			os.Exit(1)
		}
		reply = strings.Trim(reply, "\n")
		message, err := UnmarshalMessage(reply)
		if err != nil {
			fmt.Printf("Error unmarshaling server message: %s\n", err)
		} else {
			handleServerMessage(message)
		}
	}

}

func readInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	return input, err
}

func handleServerMessage(message interface{}) {
	switch message := message.(type) {
	case QuestionMessage:
		fmt.Println("Server asks us a question, that's just moronic!")
		fmt.Printf("The question was: %s\n", message.Text)
	case AnswerMessage:
		fmt.Printf("Server asnwered: %s\nAnd gave us some advice: %s\n", message.Text, message.Advice)
	default:
		log.Printf("Unsupported message type: %T", message)
	}
}
