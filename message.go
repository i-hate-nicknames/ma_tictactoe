package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	MSG_BOARD           = "board"
	MSG_WAITING_CONNECT = "waitingConnect"
	MSG_MOVE            = "move"
	SEPARATOR           = "|"
)

type WaitingMessage struct{}

type BoardMessage struct {
	board *Board
}

type MoveMessage struct {
	x, y int
}

func getMessageType(message interface{}) (string, error) {
	switch message.(type) {
	case WaitingMessage:
		return MSG_WAITING_CONNECT, nil
	case MoveMessage:
		return MSG_MOVE, nil
	case BoardMessage:
		return MSG_BOARD, nil
	default:
		return "", fmt.Errorf("%v of type %T is not a valid message to marshal", message, message)
	}
}

// MarshalMessage to a string, ready to be sent over a wire
func MarshalMessage(message interface{}) (string, error) {
	msgType, err := getMessageType(message)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return msgType + SEPARATOR + string(payload), nil
}

// UnmarshalMessage produces message of a correct type from a string
func UnmarshalMessage(marshalled string) (interface{}, error) {
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
	default:
		return nil, fmt.Errorf("Unrecognized message type %s", msgType)
	}
}
