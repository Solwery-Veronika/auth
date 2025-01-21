package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/Solwery-Veronika/auth/internal/model"
	"github.com/Solwery-Veronika/auth/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
		in := auth.RegisterRequest{
			Username: "testtest",
			Email:    "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().RegisterUser(gomock.Any(), in.Username, in.Email, in.Password).Return(model.User{Password: "123"}, nil)

		srv := New(mockRepo)
		_, err := srv.Register(ctx, &in)
		assert.ErrorContains(t, err, "invalid password")
		assert.Error(t, err)
	})

	t.Run("fail_request_error", func(t *testing.T) {
		mockErr := errors.New("mock error")
		in := auth.RegisterRequest{
			Username: "testtest",
			Email:    "testtest",
			Password: "testtest",
		}
		ctx := context.Background()
		mockRepo.EXPECT().RegisterUser(gomock.Any(), in.Username, in.Email, in.Password).Return(model.User{Password: in.Password}, mockErr)

		srv := New(mockRepo)
		_, err := srv.Register(ctx, &in)
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Error(t, err)
	})
}
