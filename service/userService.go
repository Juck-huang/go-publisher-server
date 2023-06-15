package service

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/utils"
	"time"
)

var (
	G = config.G
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

// CheckUsernameAndPassword 校验用户名和密码
func (obj *UserService) CheckUsernameAndPassword(username string, password string) error {
	rsa := utils.NewRsa(G.C.Rsa.PrivateKey)
	decrypt, err := rsa.Decrypt([]byte(password))
	if err != nil {
		G.Logger.Errorf("登录失败，失败原因:[%s]", err)
		return errors.New("用户名或密码不正确")
	}
	prepare, err := G.DB.Prepare("select count(1) from user where username = ? and password = ?")
	if err != nil {
		return err
	}
	var num int64
	err = prepare.QueryRow(username, string(decrypt)).Scan(&num)
	if err != nil {
		return err
	}
	if num > 0 {
		return nil
	}
	return errors.New("用户名或密码不正确")
}

// Logout 登出
func (obj *UserService) Logout(token string, username string) error {
	// 以用户名+过期时间+token作为key
	redisKey := fmt.Sprintf("%s_%d_%s", username, G.C.Jwt.Token.Expire, token)
	// 如果存在该key说明已经登出，直接返回已登出，如果没有则新增
	tokenRedis, err := G.RedisClient.Get(redisKey).Result()
	if err != redis.Nil {
		return err
	}
	if tokenRedis != "" {
		return errors.New("用户已登出")
	}
	err = G.RedisClient.SetNX(redisKey, token, time.Second*time.Duration(G.C.Jwt.Token.Expire)).Err()
	if err != nil {
		return err
	}
	return nil
}
