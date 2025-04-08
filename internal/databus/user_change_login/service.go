package user_change_login

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	producer *kafka.Writer
}

func New() *Service {
	producer := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "user.login.changed",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
	return &Service{
		producer: producer,
	}
}

type UserChangedLogin struct {
	Username string `json:"username"`
}

func (s *Service) SendUserChangedLogin(ctx context.Context, username string) error {
	u := UserChangedLogin{Username: username}
	uByte, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return s.producer.WriteMessages(ctx, kafka.Message{
		Value: uByte,
	})
}
