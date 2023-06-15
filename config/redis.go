package config

import (
	"fmt"
	"github.com/go-redis/redis"
)

func InitRedis() {
	G.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", G.C.Redis.Host, G.C.Redis.Port), // redis地址
		Password: G.C.Redis.Password,                                   // redis密码
		DB:       G.C.Redis.Db,                                         // redis数据库，默认0
	})
}
