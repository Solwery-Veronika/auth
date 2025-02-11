package client

import (
	"context"
	"fmt"
	"log"

	"github.com/Solwery-Veronika/auth/internal/config"
	"github.com/Solwery-Veronika/auth/internal/model"
	"github.com/Solwery-Veronika/user/pkg/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client user.UserServiceClient
}

func New(cfg *config.Config) *UserClient {
	connStr := fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port)

	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	userClient := user.NewUserServiceClient(conn)

	return &UserClient{
		client: userClient,
	}
}

func (c *UserClient) CreateUser(ctx context.Context, data model.CreateUserData) (*user.CreateUserOut, error) {
	return c.client.CreateUser(ctx, &user.CreateUserIn{
		Username: data.Username,
	})
}
