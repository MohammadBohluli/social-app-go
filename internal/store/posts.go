package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MohammadBohluli/social-app-go/types"
	"github.com/lib/pq"
)

type Post struct {
	ID        types.ID  `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    types.ID  `json:"user_id"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
	CreatedAt string    `json:"created_at"`
	Version   int       `json:"version"`
	UpdatedAt string    `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
	`

	err := s.db.
		QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s PostStore) GetByID(ctx context.Context, postID types.ID) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags, version
		FROM posts
		WHERE id = $1
	`

	var post Post
	err := s.db.QueryRowContext(ctx, query, postID).
		Scan(&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Version,
			pq.Array(&post.Tags),
		)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s PostStore) Delete(ctx context.Context, postID types.ID) error {
	query := `DELETE FROM posts WHERE id = $1`

	resp, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrorNotFound
	}

	return nil
}

func (s PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version;
	`

	err := s.db.
		QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).
		Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}
	}

	return nil
}
