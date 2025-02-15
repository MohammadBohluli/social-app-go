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
	Password  types.ID `json:"-"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (u UsersStore) Create(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := u.db.
		QueryRowContext(ctx, query, user.Username, user.Password, user.Email).
		Scan(user.ID, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
