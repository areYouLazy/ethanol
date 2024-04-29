package core

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Content   string      `json:"content"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

func (m Message) JSONMarshal() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		return nil
	}

	return msg
}

func NewMessage() Message {
	return Message{
		ID:        uuid.New().String(),
		Type:      "",
		Content:   "",
		Data:      nil,
		Timestamp: time.Now(),
	}
}

func NewErrorMessage() Message {
	msg := NewMessage()
	msg.Type = "error"

	return msg
}

func NewErrorMessageWithContent(text string) Message {
	msg := NewErrorMessage()
	msg.Content = text

	return msg
}

func NewInvalidFormatErrorMessage() Message {
	return NewErrorMessageWithContent("we speak json! you fool...")
}

func NewMarshalJSONErrorMessage() Message {
	return NewErrorMessageWithContent("internal error in json serialization while generating a response")
}

func NewIvalidTypeErrorMessage() Message {
	return NewErrorMessageWithContent("invalid message type")
}
