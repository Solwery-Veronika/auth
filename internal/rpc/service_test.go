package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/Solwery-Veronika/auth/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockDbRepo(ctrl)

	t.Run("ok", func(t *testing.T) {
		in := auth.LoginIn{
			Username: "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().LoginUser(gomock.Any(), in.Username, in.Password).Return(User{Password: in.Password}, nil)

		srv := New(mockRepo)
		_, err := srv.Login(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("fail_request_error", func(t *testing.T) {
		mockErr := errors.New("mock error")
		in := auth.LoginIn{
			Username: "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().LoginUser(gomock.Any(), in.Username, in.Password).Return(User{Password: in.Password}, mockErr)

		srv := New(mockRepo)
		_, err := srv.Login(ctx, &in)
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Error(t, err)
	})

	t.Run("fail_password_not_match", func(t *testing.T) {
		in := auth.LoginIn{
			Username: "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().LoginUser(gomock.Any(), in.Username, in.Password).Return(User{Password: "123"}, nil)

		srv := New(mockRepo)
		_, err := srv.Login(ctx, &in)
		assert.ErrorContains(t, err, "invalid password")
		assert.Error(t, err)
	})

	t.Run("fail_username_too_short", func(t *testing.T) {
		in := auth.LoginIn{
			Username: "test",
			Password: "test",
		}
		ctx := context.Background()

		srv := New(mockRepo)
		_, err := srv.Login(ctx, &in)
		assert.Error(t, err)
	})
}

func TestService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockDbRepo(ctrl)

	t.Run("ok", func(t *testing.T) {
		in := auth.SignupRequest{
			Username: "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().SignupUser(gomock.Any(), in.Username, in.Password).Return(nil)

		srv := New(mockRepo)
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

		srv := New(mockRepo)
		_, err := srv.Signup(ctx, &in)
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Error(t, err)
	})
}

func TestService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockDbRepo(ctrl)

	t.Run("ok", func(t *testing.T) {
		in := auth.RegisterUserRequest{
			Email:    "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().RegisterUser(gomock.Any(), in.Email, in.Password).Return(nil)

		srv := New(mockRepo)
		_, err := srv.RegisterUser(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("fail_request_error", func(t *testing.T) {
		mockErr := errors.New("mock error")
		in := auth.RegisterUserRequest{
			Email:    "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().RegisterUser(gomock.Any(), in.Email, in.Password).Return(mockErr)

		srv := New(mockRepo)
		_, err := srv.RegisterUser(ctx, &in)
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Error(t, err)
	})
}
