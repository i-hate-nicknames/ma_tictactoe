package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
	board        *Board
	xConn, yConn net.Conn
	serverSock   net.Listener
	numClients   int
}

func (server *Server) run() {
	for {
		// stop accepting any new clients when there are two
		if server.numClients == 2 {
			return
		}
		conn, err := server.serverSock.Accept()
		if err != nil {
			log.Printf("Failed to handle a client: %s\n", err)
			continue
		}
		if server.numClients == 0 {
			server.xConn = conn
			go server.handleClient(PLAYER_X, conn)
		} else {
			server.yConn = conn
			go server.handleClient(PLAYER_O, conn)
		}
		server.numClients++
	}
}

func startServer(port string, done chan<- bool) {
	serverSock, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start a server: %s\n", err)
	}
	done <- true
	board := MakeBoard(3)
	server := &Server{board: board, serverSock: serverSock}
	server.run()

}

// handle a client: reply to every message with modified client message
func (server *Server) handleClient(player Player, conn net.Conn) {
	reader := bufio.NewReader(conn)
	defer conn.Close()
	for {
		message, err := readMessage(reader)
		if err != nil {
			log.Printf("Error reading client message: %s\n", err)
			return
		}
		server.handleMessage(conn, message)
	}
}

// read one client message data from the given reader, parse it
// and return as a message struct
func readMessage(reader *bufio.Reader) (interface{}, error) {
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

func (server *Server) handleMessage(conn net.Conn, message interface{}) {
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
