package cache

import (
	"context"

	"github.com/MohammadBohluli/social-app-go/internal/store"
	"github.com/MohammadBohluli/social-app-go/types"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(ctx context.Context, userID types.ID) (*store.User, error)
		Set(ctx context.Context, user *store.User) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{

		Users: &UserStore{rdb: rdb},
	}
}
