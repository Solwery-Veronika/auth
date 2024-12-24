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
	user, err := s.dbR.LoginUser(ctx, in.Username, in.Password) // правильность пароля и логина
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Password != in.Password {
		return nil, status.Error(codes.Internal, "invalid password")
	}

	data := jwt.MapClaims{
		"username": in.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := token.SignedString([]byte("secret"))
	return &auth.LoginOut{
		Token: tokenString,
	}, nil // создание токена и шифрование
}

func (s *Service) Signup(ctx context.Context, in *auth.SignupRequest) (*auth.SignupResponse, error) {
	err := s.dbR.SignupUser(ctx, in.Username, in.Password) // err - ошибка от бд
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.SignupResponse{
		Success: true,
	}, nil
}

func (s *Service) RegisterUser(ctx context.Context, in *auth.RegisterUserRequest) (*auth.RegisterUserResponse, error) {
	err := s.dbR.RegisterUser(ctx, in.Email, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.RegisterUserResponse{
		Success: true,
	}, nil
}
