package rpc

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Solwery-Veronika/auth/pkg/auth"
)

type Service struct {
	auth.UnimplementedAuthServiceServer
	dbR DbRepo
}

func New(repo DbRepo) *Service {
	return &Service{
		dbR: repo,
	}
}

func (s *Service) Login(ctx context.Context, in *auth.LoginIn) (*auth.LoginOut, error) {
	err := s.dbR.LoginUser(ctx, in.Username, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := jwt.MapClaims{
		"username": in.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := token.SignedString([]byte("secret"))
	return &auth.LoginOut{
		Token: tokenString,
	}, nil
}

func (s *Service) Signup(ctx context.Context, in *auth.SignupRequest) (*auth.SignupResponse, error) {
	err := s.dbR.SignupUser(ctx, in.Username, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.SignupResponse{
		Success: true,
	}, nil
}
