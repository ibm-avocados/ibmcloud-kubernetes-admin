package notifier

import (
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Notification interface {
	Send() error
}

func GetNotifier(m *kafka.Message) (Notification, error) {
	var message Message
	if err := json.Unmarshal(m.Value, &message); err != nil {
		return nil, err
	}

	switch message.Type {
	case EMAIL:
		return NewEmail(message.Content)
	case GITHUB:
		return NewGithub(message.Content)
	default:
		return nil, errors.New("unknown message type")
	}
}

type NotificationType string

const (
	EMAIL  NotificationType = "email"
	GITHUB NotificationType = "github"
	SLACK  NotificationType = "slack"
)

type Message struct {
	Type    NotificationType `json:"type"`
	Content []byte           `json:"content"`
}
