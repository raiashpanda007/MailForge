package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type PgPool struct {
	db *pgxpool.Pool
}

func (r *PgPool) GetUserByEmail(ctx context.Context, email string) (*User, *string, error) {
	var user User
	var password string

	err := r.db.QueryRow(ctx,
		"SELECT id, name, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.Id, &user.Name, &user.Email, &password)

	if err != nil {
		return nil, nil, err
	}

	return &user, &password, nil
}

func (r *PgPool) CreateUser(ctx context.Context, email, name, password string) (*User, error) {
	var user User
	id := uuid.New()

	err := r.db.QueryRow(ctx,
		"INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, name, email",
		id, name, email, password,
	).Scan(&user.Id, &user.Name, &user.Email)

	if err != nil {
		// Unique email violation: code 23505
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, errors.New("user already exists")
		}
		return nil, err
	}

	return &user, nil
}

func (r *PgPool) DeleteUser(ctx context.Context, id string) (bool, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return false, errors.New("PLEASE PROVIDE A VALID USER ID")
	}

	cmd, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return false, err
	}

	if cmd.RowsAffected() == 0 {
		return false, errors.New("USER DOESN'T EXIST OF THIS ID")
	}

	return true, nil
}

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*User, *string, error)
	CreateUser(ctx context.Context, email string, name string, password string) (*User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
}

func NewUserRepo(db *pgxpool.Pool) UserRepo {
	return &PgPool{db}
}
