package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/Solwery-Veronika/auth/internal/config"
	"github.com/Solwery-Veronika/auth/internal/model"
	"github.com/Solwery-Veronika/auth/pkg/auth"
	"github.com/Solwery-Veronika/user/pkg/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockDbRepo(ctrl)
	mockUser := NewMockUserC(ctrl)
	cfg := &config.Config{}

	t.Run("ok", func(t *testing.T) {
		in := auth.SignupRequest{
			Username: "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().SignupUser(gomock.Any(), in.Username, in.Password).Return(nil)
		mockUser.EXPECT().CreateUser(gomock.Any(), model.CreateUserData{Username: in.Username}).Return(&user.CreateUserOut{Success: true}, nil)
		srv := New(cfg, mockRepo, mockUser)
		_, err := srv.Signup(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("fail_request_error", func(t *testing.T) {
		mockErr := errors.New("mock error")
		in := auth.SignupRequest{
			Username: "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().SignupUser(gomock.Any(), in.Username, in.Password).Return(mockErr)

		srv := New(cfg, mockRepo, mockUser)
		_, err := srv.Signup(ctx, &in)
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Error(t, err)
	})
}

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockDbRepo(ctrl)
	mockUser := NewMockUserC(ctrl)
	cfg := &config.Config{}

	t.Run("ok", func(t *testing.T) {
		in := auth.LoginIn{
			Username: "testtest",
			Email:    "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().LoginUser(gomock.Any(), in.Username, in.Email, in.Password).Return(model.User{Password: "123"}, nil)

		srv := New(cfg, mockRepo, mockUser)
		_, err := srv.Login(ctx, &in)
		assert.ErrorContains(t, err, "invalid password")
		assert.Error(t, err)
	})

	t.Run("fail_request_error", func(t *testing.T) {
		mockErr := errors.New("mock error")
		in := auth.LoginIn{
			Username: "testtest",
			Email:    "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().LoginUser(gomock.Any(), in.Username, in.Email, in.Password).Return(model.User{Password: in.Password}, mockErr)

		srv := New(cfg, mockRepo, mockUser)
		_, err := srv.Login(ctx, &in)
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Error(t, err)
	})
}

func TestService_Change(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockDbRepo(ctrl)
	mockUser := NewMockUserC(ctrl)
	cfg := &config.Config{}

	t.Run("ok", func(t *testing.T) {
		in := auth.ChangeLoginIn{
			Username:    "testtest",
			Password:    "testtest",
			NewUsername: "newUser",
		}
		ctx := context.Background()
		mockRepo.EXPECT().LoginUser(gomock.Any(), in.Username, "", in.Password).Return(model.User{}, nil)

		mockRepo.EXPECT().LoginUser(gomock.Any(), in.NewUsername, "", "").Return(model.User{}, errors.New("user not found"))

		mockRepo.EXPECT().ChangeLogin(gomock.Any(), in.Username, in.Password, in.NewUsername).Return(model.ChangeUser{NewUsername: in.NewUsername}, nil)

		srv := New(cfg, mockRepo, mockUser)
		_, err := srv.ChangeLogin(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("fail_request_error", func(t *testing.T) {
		mockErr := errors.New("mock error")
		in := auth.ChangeLoginIn{
			Username:    "testtest",
			Password:    "testtest",
			NewUsername: "newUser",
		}
		ctx := context.Background()
		// 1. Проверяем старый логин (возвращает успешный вход)
		mockRepo.EXPECT().LoginUser(gomock.Any(), in.Username, "", in.Password).Return(model.User{}, nil)

		// 2. Проверяем, что новый логин не занят (возвращает ошибку, значит, логин свободен)
		mockRepo.EXPECT().LoginUser(gomock.Any(), in.NewUsername, "", "").Return(model.User{}, errors.New("user not found"))

		mockRepo.EXPECT().
			ChangeLogin(gomock.Any(), in.Username, in.Password, in.NewUsername).
			Return(model.ChangeUser{}, mockErr)

		srv := New(cfg, mockRepo, mockUser)
		_, err := srv.ChangeLogin(ctx, &in)

		// Проверяем, что ошибка присутствует
		assert.Error(t, err)
		assert.ErrorContains(t, err, mockErr.Error()) // Проверяем, что текст ошибки совпадает
	})
}
