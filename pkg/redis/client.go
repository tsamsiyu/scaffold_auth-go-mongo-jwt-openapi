package redis

import (
	"github.com/go-redis/redis/v8"
)

func NewClient(uri string, dbnum int, pass string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: pass,
		DB:       dbnum,
	})
}
