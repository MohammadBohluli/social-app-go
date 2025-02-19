package store

import (
	"context"
	"database/sql"
	"errors"

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
	}
	Comments interface {
		Create(context.Context, Comment) error
		GetByPostID(ctx context.Context, postID types.ID) ([]Comment, error)
	}

	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, types.ID) (*User, error)
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
