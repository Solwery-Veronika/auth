package rpc

import (
	"context"
	"errors"
	"log"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Solwery-Veronika/auth/internal/config"
	"github.com/Solwery-Veronika/auth/internal/model"
	"github.com/Solwery-Veronika/auth/internal/repository/postgres"
	"github.com/Solwery-Veronika/auth/pkg/auth"
)

type Service struct {
	auth.UnimplementedAuthServiceServer
	dbR         DbRepo
	secretToken string
	uC          UserC
}

func New(cfg *config.Config, repo DbRepo, uC UserC) *Service {
	return &Service{
		dbR:         repo,
		secretToken: cfg.Platform.Secret,
		uC:          uC,
	}
}

func (s *Service) Login(ctx context.Context, in *auth.LoginIn) (*auth.LoginOut, error) {
	if len(in.Username) < 8 {
		return nil, status.Error(codes.InvalidArgument, "username too short")
	}
	user, err := s.dbR.LoginUser(ctx, in.Username, in.Email, in.Password) // правильность пароля и логина
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Password != in.Password {
		return nil, status.Error(codes.Internal, "invalid password")
	}

	data := jwt.MapClaims{
		"username": in.Username,
		"email":    in.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := token.SignedString([]byte(s.secretToken))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.LoginOut{
		Token: tokenString,
	}, nil // создание токена и шифрование
}

func (s *Service) Signup(ctx context.Context, in *auth.SignupRequest) (*auth.SignupResponse, error) {
	log.Println("signup work")
	err := s.dbR.SignupUser(ctx, in.Username, in.Password) // err - ошибка от бд
	success := true
	if err != nil {
		if !errors.Is(err, postgres.ErrUserExists) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		success = false
	}

	res, err := s.uC.CreateUser(ctx, model.CreateUserData{
		Username: in.Username,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.SignupResponse{
		Success: success && res.Success,
	}, nil
}
