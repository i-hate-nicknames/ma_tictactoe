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
		input, err := readInput(inputReader)
		input = strings.Trim(input, "\n")
		if err == io.EOF {
			fmt.Println("Exiting")
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("Error reading input: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Sending " + input)
		fmt.Fprintln(conn, input)
		reply, err := sockReader.ReadString('\n')
		reply = strings.Trim(reply, "\n")
		if err == io.EOF {
			fmt.Println("Server closed the connection")
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("Error reading reply from server: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Server replied: " + reply)
	}

}

func readInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	return input, err
}
