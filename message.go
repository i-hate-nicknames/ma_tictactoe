package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	MSG_QUESTION = "QUE"
	MSG_ANSWER   = "ANS"
	SEPARATOR    = "|"
)

type QuestionMessage struct {
	Text string
}

type AnswerMessage struct {
	Text   string
	Advice string
}

func getMessageType(message interface{}) (string, error) {
	switch message.(type) {
	case QuestionMessage:
		return MSG_QUESTION, nil
	case AnswerMessage:
		return MSG_ANSWER, nil
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
	case MSG_QUESTION:
		var message QuestionMessage
		err := json.Unmarshal(payload, &message)
		if err != nil {
			return nil, err
		}
		return message, nil
	case MSG_ANSWER:
		var message AnswerMessage
		err := json.Unmarshal(payload, &message)
		if err != nil {
			return nil, err
		}
		return message, nil
	default:
		return nil, fmt.Errorf("Unrecognized message type %s", msgType)
	}
}
