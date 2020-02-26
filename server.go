package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	board        *Board
	gameLock     sync.Mutex
	xConn, yConn net.Conn
	serverSock   net.Listener
	numClients   int
	gameStarted  bool
}

func (server *Server) run() {
	for {
		// stop accepting any new clients when the game has started
		if server.gameStarted {
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
		if server.numClients == 2 {
			server.gameStarted = true
		}
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
	sendMessage(conn, HelloMessage{"Welcome to this tic tac toe server!", player})
	if server.gameStarted {
		sendMessage(conn, BoardMessage{server.board})
	} else {
		sendMessage(conn, WaitingMessage{})
		for !server.gameStarted {
			time.Sleep(300 * time.Millisecond)
		}
		sendMessage(conn, BoardMessage{server.board})
	}
	for {
		message, err := readMessage(reader)
		if err != nil {
			log.Printf("Error reading client message: %s\n", err)
			return
		}
		server.handleMessage(player, conn, message)
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

func (server *Server) handleMessage(player Player, conn net.Conn, message interface{}) {
	if !server.gameStarted {
		// ignore client messages until the game has started
		return
	}
	switch message := message.(type) {
	case MoveMessage:
		server.gameLock.Lock()
		err := server.board.MakeMove(player, message.X, message.Y)
		server.gameLock.Unlock()
		if err != nil {
			sendMessage(conn, ErrorMessage{err.Error()})
		} else {
			sendMessage(conn, BoardMessage{server.board})
		}
	default:
		log.Printf("Unsupported message type: %T", message)
	}
}

// todo check if this is a blocking call
func sendMessage(conn net.Conn, message interface{}) {
	data, err := MarshalMessage(message)
	if err != nil {
		log.Printf("Error marshaling message: %s\n", err)
		return
	}
	fmt.Fprintln(conn, data)
}
