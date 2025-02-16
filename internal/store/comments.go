package store

import (
	"context"
	"database/sql"

	"github.com/MohammadBohluli/social-app-go/types"
)

type Comment struct {
	ID        types.ID `json:"id"`
	PostID    types.ID `json:"post_id"`
	UserID    types.ID `json:"user_id"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	User      User     `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s CommentStore) GetByPostID(ctx context.Context, postID types.ID) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id FROM comments c
		JOIN users ON users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;
	`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return []Comment{}, err
	}

	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		)
		if err != nil {
			return []Comment{}, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (c CommentStore) Create(ctx context.Context, comment Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	err := c.db.
		QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
