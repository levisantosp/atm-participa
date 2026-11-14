package redis

import (
	"github.com/redis/go-redis/v9"
)

const Nil = redis.Nil

var Client *redis.Client

func Connect(addr, password string) {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
}
