package local

import (
	"context"
	"fmt"
)

type Repository struct {
	storage map[string]string
}

func NewRepository() *Repository {
	storage := make(map[string]string)
	return &Repository{storage: storage}
}

func (r *Repository) SignupUser(ctx context.Context, username string, password string) error {
	_, ok := r.storage[username]
	if ok {
		return fmt.Errorf("username %s already exists", username)
	}
	r.storage[username] = password
	return nil
}

func (r *Repository) LoginUser(ctx context.Context, username string, password string) error {
	pass, ok := r.storage[username]
	if !ok {
		return fmt.Errorf("username or password not match")
	}
	if pass != password {
		return fmt.Errorf("username or password not match")
	}
	return nil
}
