package service

import (
	"errors"
	"fmt"
	"hy.juck.com/go-publisher-server/model"
	"time"

	"github.com/go-redis/redis"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/utils"
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
	rsa := utils.NewRsa(G.C.Jwt.Rsa.PrivateKey)
	decryptPassword, err := rsa.Decrypt([]byte(password))
	if err != nil {
		G.Logger.Errorf("解析密码失败，失败原因:[%s]，加密密码:[%s]", err, password)
		return errors.New("用户名或密码不正确")
	}
	// 校验现有的密码和argon2生成的加密密码正确性
	var user model.User
	G.DB.Debug().Where("username = ?", username).First(&user)
	if user.Password == "" {
		G.Logger.Errorf("用户名[%s]不存在", username)
		return errors.New("用户名不存在")
	}
	// 1.先从数据库查询出现有用户名对应加密后的密码
	p := &utils.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	argonUtils := utils.NewArgon2(p)
	// 2.数据库取出来的argon2加密后的密码和现有的做比对
	matchPassword, _ := argonUtils.ComparePasswordAndHash(string(decryptPassword), user.Password)
	if matchPassword {
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
