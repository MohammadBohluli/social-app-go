package cache

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func New(host string, port int, password string, db int) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	log.Println("✅ Redis is up running...")

	return rdb
}
