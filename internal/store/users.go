package store

import (
	"context"
	"database/sql"

	"github.com/MohammadBohluli/social-app-go/types"
)

type User struct {
	ID        types.ID `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  string   `json:"-"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type UserStore struct {
	db *sql.DB
}

func (u UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := u.db.
		QueryRowContext(ctx, query, user.Username, user.Password, user.Email).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s UserStore) GetByID(ctx context.Context, userID types.ID) (*User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := s.db.QueryRowContext(ctx, query, userID).
		Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
