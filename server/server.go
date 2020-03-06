package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"nvm.ga/mastersofcode/golang_2019/tictactoe/game"
	msg "nvm.ga/mastersofcode/golang_2019/tictactoe/messaging"
)

type Server struct {
	board       *game.Board
	gameLock    sync.Mutex
	serverSock  net.Listener
	numClients  int
	gameStarted bool
}

func StartServer(port string, done chan<- bool) {
	serverSock, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start a server: %s\n", err)
	}
	done <- true
	board := game.MakeBoard(3)
	server := &Server{board: board, serverSock: serverSock}
	server.run()
}

type ConnectedPlayer struct {
	player          game.Player
	conn            net.Conn
	opponentUpdates <-chan bool
	myUpdates       chan<- bool
}

func (server *Server) run() {
	updatesX := make(chan bool)
	updatesO := make(chan bool)
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
			connPlayer := &ConnectedPlayer{
				player:          game.PLAYER_X,
				conn:            conn,
				opponentUpdates: updatesO,
				myUpdates:       updatesX,
			}
			go server.handleClient(connPlayer)
		} else {
			connPlayer := &ConnectedPlayer{
				player:          game.PLAYER_O,
				conn:            conn,
				opponentUpdates: updatesX,
				myUpdates:       updatesO,
			}
			go server.handleClient(connPlayer)
		}

		server.numClients++
		if server.numClients == 2 {
			server.gameStarted = true
		}
	}
}

// handle a client: reply to every message with modified client message
func (server *Server) handleClient(connPlayer *ConnectedPlayer) {
	defer connPlayer.conn.Close()
	connPlayer.sendMessage(msg.HelloMessage{"Welcome to this tic tac toe server!", connPlayer.player})
	if server.gameStarted {
		connPlayer.sendBoard(server.board)
	} else {
		connPlayer.sendMessage(msg.WaitingMessage{})
		for !server.gameStarted {
			time.Sleep(300 * time.Millisecond)
		}
		connPlayer.sendBoard(server.board)
	}
	clientChan := make(chan msg.Message, 0)
	errChan := make(chan error)
	go msg.ReadMessages(connPlayer.conn, clientChan, errChan)
	for {
		select {
		case clientMessage := <-clientChan:
			server.handleMessage(connPlayer, clientMessage)
		case <-connPlayer.opponentUpdates:
			// todo: read a string from opponent updates, and dispatch on it
			// handle disconnected opponent gracefuly (add Exit Message)
			connPlayer.sendBoard(server.board)
		case err := <-errChan:
			if err == io.EOF {
				fmt.Println("client disconnected")
				// todo: maybe send update to the other client
				os.Exit(1)
			} else {
				fmt.Println("Error reading message from client " + err.Error())
				// todo: maybe check the error and ignore if it's not fatal,
				// i.e. malformed message
				os.Exit(1)
			}
		}
	}
}

func (connPlayer *ConnectedPlayer) sendMessage(message msg.Message) {
	msg.SendMessage(connPlayer.conn, message)
}

func (connPlayer *ConnectedPlayer) sendBoard(board *game.Board) {
	msg.SendMessage(connPlayer.conn, msg.BoardMessage{board})
}

func (server *Server) handleMessage(connPlayer *ConnectedPlayer, message msg.Message) {
	if !server.gameStarted {
		// ignore client messages until the game has started
		return
	}
	switch message := message.(type) {
	case msg.MoveMessage:
		server.gameLock.Lock()
		err := server.board.MakeMove(connPlayer.player, message.X, message.Y)
		server.gameLock.Unlock()
		if err != nil {
			msg.SendMessage(connPlayer.conn, msg.ErrorMessage{err.Error()})
		} else {
			msg.SendMessage(connPlayer.conn, msg.BoardMessage{server.board})
			connPlayer.myUpdates <- true
		}
	default:
		log.Printf("Unsupported message type: %T", message)
	}
}
