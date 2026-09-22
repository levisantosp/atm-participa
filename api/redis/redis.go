package redis

import (
	"github.com/levisantosp/atm-participa/api/utils"
	"github.com/redis/go-redis/v9"
)

const Nil = redis.Nil

var Client *redis.Client

func Connect() {
	Client = redis.NewClient(&redis.Options{
		Addr:     utils.Env.RedisAddr,
		Password: utils.Env.RedisPassword,
	})
}
