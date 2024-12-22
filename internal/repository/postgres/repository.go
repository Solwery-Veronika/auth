package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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
	query := `INSERT INTO participants (username, password) VALUES ($1, $2)` // запрос

	_, err := r.conn.ExecContext(ctx, query, username, password)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %w", err)
	}
	return nil
}

type User struct {
	Password string `db:"password"`
}

func (r *Repository) LoginUser(ctx context.Context, username string, password string) error {
	query := `SELECT password FROM participants WHERE username = $1`

	var user User

	err := r.conn.Get(&user, query, username) // заполнение структуры результатом запроса
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Password != password {
		return fmt.Errorf("invalid password")
	}

	return nil
}
