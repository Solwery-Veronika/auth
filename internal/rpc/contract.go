package rpc

import "context"

type DbRepo interface {
	SignupUser(ctx context.Context, username string, password string) error
	LoginUser(ctx context.Context, username string, password string) error
}
