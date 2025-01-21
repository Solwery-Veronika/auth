//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package rpc

import (
	"context"

	"github.com/Solwery-Veronika/auth/internal/model"
)

type DbRepo interface {
	SignupUser(ctx context.Context, username string, password string) error
	// LoginUser(ctx context.Context, username string, password string) (model.User, error)
	// RegisterUser(ctx context.Context, email string, password string) error
	LoginUser(ctx context.Context, username string, email string, password string) (model.User, error)
}
