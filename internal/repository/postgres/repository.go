package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Solwery-Veronika/auth/internal/config"
	"github.com/Solwery-Veronika/auth/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrUserExists = errors.New("user already exists")

type Repository struct {
	conn *sqlx.DB
}

// Проверяем, существует ли пользователь с таким логином и паролем
func (r *Repository) getUserID(ctx context.Context, username, password string) (int, error) {
	var userID int

	query := "SELECT id FROM participants WHERE username = $1 AND password = $2"

	err := r.conn.QueryRowContext(ctx, query, username, password).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("invalid username or password")
		}

		return 0, err
	}

	return userID, nil
}

// Проверяем, что новый логин свободен
func (r *Repository) isUsernameTaken(ctx context.Context, username string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM participants WHERE username = $1"
	err := r.conn.QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Обновляем
func (r *Repository) updateUsername(ctx context.Context, userID int, newUsername string) error {
	query := "UPDATE participants SET username = $1 WHERE id = $2"
	_, err := r.conn.ExecContext(ctx, query, newUsername, userID)
	return err
}

///////////////////////////

func NewRepository(cfg *config.Config) *Repository {
	connectCmd := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port) // строка для подключения к pg

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
	query := `SELECT true FROM participants WHERE username = $1`
	var user user
	var exists bool

	err := r.conn.GetContext(ctx, &exists, query, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	if exists {
		return model.User{}, fmt.Errorf("user already exists")
	}

	// Добавляем нового пользователя в базу данных
	queryInsert := `INSERT INTO participants (username, email, password) VALUES ($1, $2, $3);`
	_, err = r.conn.ExecContext(ctx, queryInsert, username, email, password)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return model.User{
		Password: user.Password,
	}, nil
}

func (r *Repository) ChangeLogin(ctx context.Context, username, password, newUsername string) (model.ChangeUser, error) {
	userID, err := r.getUserID(ctx, username, password)
	if err != nil {
		return model.ChangeUser{}, err
	}

	taken, err := r.isUsernameTaken(ctx, newUsername)
	if err != nil {
		return model.ChangeUser{}, err
	}
	if taken {
		return model.ChangeUser{}, errors.New("username already taken")
	}

	err = r.updateUsername(ctx, userID, newUsername)
	if err != nil {
		return model.ChangeUser{}, err
	}

	return model.ChangeUser{NewUsername: newUsername}, nil
}
