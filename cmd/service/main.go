package main

import (
	"fmt"
	"github.com/Solwery-Veronika/auth/internal/databus/user_change_login"
	"log"
	"net"

	"github.com/Solwery-Veronika/auth/internal/client"
	"github.com/Solwery-Veronika/auth/internal/config"
	"github.com/Solwery-Veronika/auth/internal/repository/postgres"
	"github.com/Solwery-Veronika/auth/internal/rpc"
	"github.com/Solwery-Veronika/auth/pkg/auth"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	repo := postgres.NewRepository(cfg)
	kafkaProducer := user_change_login.New()
	//err := kafkaProducer.SendUserChangedLogin(context.Background(), "test 2")
	//if err != nil {
	//	log.Fatal(err)
	//}

	userClient := client.New(cfg)

	service := rpc.New(cfg, repo, userClient, kafkaProducer)

	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, service)

	log.Println(cfg)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
