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
	xConn, yConn *ConnectedPlayer
	serverSock   net.Listener
	numClients   int
	gameStarted  bool
}

type ConnectedPlayer struct {
	player       Player
	conn         net.Conn
	stateUpdates chan bool
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
		stateUpdates := make(chan bool)
		if server.numClients == 0 {
			connPlayer := &ConnectedPlayer{player: PLAYER_X, conn: conn, stateUpdates: stateUpdates}
			server.xConn = connPlayer
			go server.handleClient(connPlayer)
		} else {
			connPlayer := &ConnectedPlayer{player: PLAYER_O, conn: conn, stateUpdates: stateUpdates}
			server.yConn = connPlayer
			go server.handleClient(connPlayer)
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
func (server *Server) handleClient(connPlayer *ConnectedPlayer) {
	reader := bufio.NewReader(connPlayer.conn)
	defer connPlayer.conn.Close()
	sendMessage(connPlayer, HelloMessage{"Welcome to this tic tac toe server!", connPlayer.player})
	if server.gameStarted {
		sendMessage(connPlayer, BoardMessage{server.board})
	} else {
		sendMessage(connPlayer, WaitingMessage{})
		for !server.gameStarted {
			time.Sleep(300 * time.Millisecond)
		}
		sendMessage(connPlayer, BoardMessage{server.board})
	}
	for {
		message, err := readMessage(reader)
		if err != nil {
			log.Printf("Error reading client message: %s\n", err)
			return
		}
		server.handleMessage(connPlayer, message)
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

func (server *Server) handleMessage(connPlayer *ConnectedPlayer, message interface{}) {
	if !server.gameStarted {
		// ignore client messages until the game has started
		return
	}
	switch message := message.(type) {
	case MoveMessage:
		server.gameLock.Lock()
		err := server.board.MakeMove(connPlayer.player, message.X, message.Y)
		server.gameLock.Unlock()
		if err != nil {
			sendMessage(connPlayer, ErrorMessage{err.Error()})
		} else {
			sendMessage(connPlayer, BoardMessage{server.board})
		}
	default:
		log.Printf("Unsupported message type: %T", message)
	}
}

// todo check if this is a blocking call
func sendMessage(connPlayer *ConnectedPlayer, message interface{}) {
	data, err := MarshalMessage(message)
	if err != nil {
		log.Printf("Error marshaling message: %s\n", err)
		return
	}
	fmt.Fprintln(connPlayer.conn, data)
}
