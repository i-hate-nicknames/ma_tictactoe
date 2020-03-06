package messaging

import (
	"encoding/json"
	"fmt"
	"strings"

	"nvm.ga/mastersofcode/golang_2019/tictactoe/game"
)

const (
	MSG_BOARD           = "board"
	MSG_WAITING_CONNECT = "waitingConnect"
	MSG_MOVE            = "move"
	MSG_ERROR           = "error"
	MSG_HELLO           = "hello"
	SEPARATOR           = "|"
)

type WaitingMessage struct{}

type BoardMessage struct {
	Board *game.Board
}

type MoveMessage struct {
	X, Y int
}

type ErrorMessage struct {
	Text string
}

type HelloMessage struct {
	Text           string
	AssignedPlayer game.Player
}

type Message interface {
	GetType() string
}

func (m HelloMessage) GetType() string {
	return MSG_HELLO
}

func (m WaitingMessage) GetType() string {
	return MSG_WAITING_CONNECT
}

func (m MoveMessage) GetType() string {
	return MSG_MOVE
}

func (m BoardMessage) GetType() string {
	return MSG_BOARD
}

func (m ErrorMessage) GetType() string {
	return MSG_ERROR
}

// MarshalMessage to a string, ready to be sent over a wire
func MarshalMessage(message Message) (string, error) {
	payload, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return message.GetType() + SEPARATOR + string(payload), nil
}

// UnmarshalMessage produces message of a correct type from a string
func UnmarshalMessage(marshalled string) (Message, error) {
	parts := strings.Split(marshalled, SEPARATOR)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Unmarshal: malformed message: %s", marshalled)
	}
	msgType, payload := parts[0], []byte(parts[1])
	switch msgType {
	case MSG_BOARD:
		var message BoardMessage
		err := json.Unmarshal(payload, &message)
		if err != nil {
			return nil, err
		}
		return message, nil
	case MSG_MOVE:
		var message MoveMessage
		err := json.Unmarshal(payload, &message)
		if err != nil {
			return nil, err
		}
		return message, nil
	case MSG_WAITING_CONNECT:
		var message WaitingMessage
		err := json.Unmarshal(payload, &message)
		if err != nil {
			return nil, err
		}
		return message, nil
	case MSG_ERROR:
		var message ErrorMessage
		err := json.Unmarshal(payload, &message)
		if err != nil {
			return nil, err
		}
		return message, nil
	case MSG_HELLO:
		var message HelloMessage
		err := json.Unmarshal(payload, &message)
		if err != nil {
			return nil, err
		}
		return message, nil
	default:
		return nil, fmt.Errorf("Unrecognized message type %s", msgType)
	}
}
