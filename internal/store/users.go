package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MohammadBohluli/social-app-go/pkg"
	"github.com/MohammadBohluli/social-app-go/types"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type User struct {
	ID        types.ID `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  string   `json:"-"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	IsActive  bool     `json:"is_active"`
}

type UserStore struct {
	db *sql.DB
}

func (u UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	hashPassword, _ := pkg.Hash(user.Password)
	err := u.db.
		QueryRowContext(ctx, query, user.Username, hashPassword, user.Email).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}

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

func (s UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTX(s.db, ctx, func(tx *sql.Tx) error {

		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createAndInvite(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})

}

func (s UserStore) createAndInvite(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID types.ID) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3);`

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}

func (s UserStore) Activate(ctx context.Context, token string) error {
	return withTX(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
    	SELECT u.id, u.username, u.email, u.created_at, u.is_active
    	FROM users u
    	JOIN user_invitations ui ON u.id = ui.user_id
    	WHERE ui.token = $1 AND ui.expiry > $2
	`

	user := &User{}
	err := tx.QueryRowContext(ctx, query, token, time.Now()).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := "UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4"

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID types.ID) error {
	query := "DELETE FROM user_invitations WHERE user_id = $1"

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s UserStore) Delete(ctx context.Context, userID types.ID) error {
	return withTX(s.db, ctx, func(tx *sql.Tx) error {

		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userID types.ID) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at FROM users
		WHERE email = $1 AND is_active = true;
	`

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
