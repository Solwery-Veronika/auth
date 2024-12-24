package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Solwery-Veronika/auth/internal/rpc"
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

type user struct {
	Password string `db:"password"`
}

func (r *Repository) LoginUser(ctx context.Context, username string, password string) (rpc.User, error) {
	query := `SELECT password FROM participants WHERE username = $1`

	var user user

	err := r.conn.Get(&user, query, username) // заполнение структуры результатом запроса
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rpc.User{}, fmt.Errorf("user not found")
		}

		return rpc.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return rpc.User{
		Password: user.Password,
	}, nil
}

func (r *Repository) RegisterUser(ctx context.Context, email string, password string) error {
	query := `SELECT true FROM participants WHERE email = $1`

	var exists bool

	err := r.conn.GetContext(ctx, &exists, query, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if exists {
		return fmt.Errorf("user already exists")
	}

	// Добавляем нового пользователя в базу данных
	queryInsert := `INSERT INTO participants (email, password) VALUES ($1, $2);`
	_, err = r.conn.ExecContext(ctx, queryInsert, email, password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
