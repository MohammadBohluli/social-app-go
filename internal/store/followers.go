package store

import (
	"context"
	"database/sql"

	"github.com/MohammadBohluli/social-app-go/types"
	"github.com/lib/pq"
)

type Follower struct {
	UserID     types.ID `json:"user_id"`
	FollowerID types.ID `json:"follower_id"`
	CreatedAt  string   `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (f FollowerStore) Follow(ctx context.Context, followerID, userID types.ID) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2);
	`

	_, err := f.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		if psqlErr, ok := err.(*pq.Error); ok && psqlErr.Code == "23505" {
			return ErrorConflict
		}
		return err
	}

	return nil
}
func (f FollowerStore) UnFollow(ctx context.Context, followerID, userID types.ID) error {

	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2;
	`

	_, err := f.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	return nil
}
