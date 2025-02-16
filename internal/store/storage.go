package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MohammadBohluli/social-app-go/types"
)

var (
	ErrorNotFound = errors.New("Resource not found")
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
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    PostStore{db},
		Users:    UserStore{db},
		Comments: CommentStore{db},
	}
}
