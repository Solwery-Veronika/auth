//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package rpc

import (
	"context"

	"github.com/Solwery-Veronika/auth/internal/model"
	"github.com/Solwery-Veronika/user/pkg/user"
)

type DbRepo interface {
	SignupUser(ctx context.Context, username string, password string) error
	LoginUser(ctx context.Context, username string, email string, password string) (model.User, error)
}
type UserC interface {
	CreateUser(ctx context.Context, data model.CreateUserData) (*user.CreateUserOut, error)
}
