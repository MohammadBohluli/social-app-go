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
	ErrorNotFound = errors.New("resource not found")
	ErrorConflict = errors.New("resource already exists")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		Update(context.Context, *Post) error
		Delete(context.Context, types.ID) error
		GetByID(context.Context, types.ID) (*Post, error)
		GetUserFeed(context.Context, types.ID, pkg.PaginationFeedQuery) ([]PostWithMetaData, error)
	}
	Comments interface {
		Create(context.Context, Comment) error
		GetByPostID(ctx context.Context, postID types.ID) ([]Comment, error)
	}

	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		Activate(context.Context, string) error
		GetByID(context.Context, types.ID) (*User, error)
		Delete(context.Context, types.ID) error
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
	}

	Followers interface {
		Follow(ctx context.Context, followerID, userID types.ID) error
		UnFollow(ctx context.Context, followerID, userID types.ID) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     PostStore{db},
		Users:     UserStore{db},
		Comments:  CommentStore{db},
		Followers: FollowerStore{db},
	}
}

func withTX(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
