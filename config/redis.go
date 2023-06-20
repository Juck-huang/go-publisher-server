package config

import (
	"fmt"

	"github.com/go-redis/redis"
)

func InitRedis() {
	G.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", G.C.Application.Redis.Host, G.C.Application.Redis.Port), // redis地址
		Password: G.C.Application.Redis.Password,                                               // redis密码
		DB:       G.C.Application.Redis.Db,                                                     // redis数据库，默认0
	})
}
