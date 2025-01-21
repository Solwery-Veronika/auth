package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Solwery-Veronika/auth/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrUserExists = errors.New("user already exists")

type Repository struct {
	conn *sqlx.DB
}

func NewRepository() *Repository {
	connectCmd := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		"master", "master", "master", "localhost", "3115") // строка для подключения к pg

	conn, err := sqlx.Connect("postgres", connectCmd) // подключаемся к бд
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{conn: conn}
}

func (r *Repository) SignupUser(ctx context.Context, username string, password string) error {
	query := `SELECT true FROM participants WHERE username = $1`

	var exists bool

	err := r.conn.GetContext(ctx, &exists, query, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if exists {
		return ErrUserExists
	}

	query = `INSERT INTO participants (username, password) VALUES ($1, $2)` // запрос

	_, err = r.conn.ExecContext(ctx, query, username, password)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %w", err)
	}
	return nil
}

type user struct {
	Password string `db:"password"`
}

func (r *Repository) LoginUser(ctx context.Context, username string, email string, password string) (model.User, error) {
	query := `SELECT true FROM participants WHERE = $1`
	var user user
	var exists bool

	err := r.conn.GetContext(ctx, &exists, query, username, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	if exists {
		return model.User{}, fmt.Errorf("user already exists")
	}

	// Добавляем нового пользователя в базу данных
	queryInsert := `INSERT INTO participants (username, email, password) VALUES ($1, $2);`
	_, err = r.conn.ExecContext(ctx, queryInsert, username, email, password)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return model.User{
		Password: user.Password,
	}, nil
}
