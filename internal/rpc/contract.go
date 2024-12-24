//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package rpc

import "context"

type User struct {
	Password string
}

type DbRepo interface {
	SignupUser(ctx context.Context, username string, password string) error
	LoginUser(ctx context.Context, username string, password string) (User, error)
	RegisterUser(ctx context.Context, email string, password string) error
}
