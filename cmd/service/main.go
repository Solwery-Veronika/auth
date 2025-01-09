package main

import (
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/Solwery-Veronika/auth/internal/repository/postgres"
	"github.com/Solwery-Veronika/auth/internal/rpc"
	"github.com/Solwery-Veronika/auth/pkg/auth"
)

func main() {
	repo := postgres.NewRepository()

	service := rpc.New(repo)

	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, service)

	lis, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
