package rpc

import (
	"context"
	"errors"
	"github.com/Solwery-Veronika/auth/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
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
