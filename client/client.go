package client

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"nvm.ga/mastersofcode/golang_2019/tictactoe/game"
	msg "nvm.ga/mastersofcode/golang_2019/tictactoe/messaging"
)

type Client struct {
	conn       net.Conn
	connReader *bufio.Reader
	player     game.Player
	board      *game.Board
}

func StartClient(addr string) {
	log.Println("Starting client, connecting to " + addr)
	conn, err := net.Dial("tcp4", addr)
	defer conn.Close()
	if err != nil {
		log.Fatalf("Couldn't connect to %s", addr)
		return
	}
	sockReader := bufio.NewReader(conn)
	client := &Client{conn: conn, connReader: sockReader, player: game.NO_PLAYER}
	for {
		message := msg.ReadServerMessage(sockReader)
		client.handleMessage(message)
	}
}

func (client *Client) handleMessage(message interface{}) {
	switch message := message.(type) {
	case msg.WaitingMessage:
		fmt.Println("Waiting for another player to connect")
	case msg.BoardMessage:
		board := message.Board
		client.board = board
		// print state
		fmt.Printf("Server sent us a board!\n%s\n", board)
		if board.GetState() != game.PLAYING {
			fmt.Println("Game over")
			return
		}
		if client.player != board.NextTurn {
			fmt.Println("Waiting for the opponent")
			return
		}
		reply, err := readMove(board)
		if err != nil {
			fmt.Println(err)
			// retry handling
			client.handleMessage(message)
		}
		msg.SendMessage(client.conn, reply)
	case msg.ErrorMessage:
		fmt.Printf("Error: %s\n", message.Text)
		if client.board == nil {
			return
		}
		// todo: Add multiple error messages, or error type to ErrorMessage
		// dispatch on that and react accordingly: ask client to repeat input if
		// it was invalid, or just show the message
		// assuming the Error was an incorrect move, retry reading
		// user input
		reply, err := readMove(client.board)
		if err != nil {
			fmt.Println(err)
			// retry handling
			client.handleMessage(message)
		}
		msg.SendMessage(client.conn, reply)
	case msg.HelloMessage:
		fmt.Printf("%s\nYour player is %s\n", message.Text, message.AssignedPlayer)
		client.player = message.AssignedPlayer
	default:
		log.Printf("Unsupported message type: %T", message)
	}
}

func readMove(board *game.Board) (msg.MoveMessage, error) {
	var result msg.MoveMessage
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
	err = board.ValidateCoordinates(x, y)
	if err != nil {
		return result, err
	}
	return msg.MoveMessage{x, y}, nil
}
