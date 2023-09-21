package service

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/dto/user"
	user2 "hy.juck.com/go-publisher-server/model/user"
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

// GetUserInfo 通过token获取当前登录用户信息
func (obj *UserService) GetUserInfo(username string) (userDto user.ResponseDto, err error) {
	var userModel user2.User
	G.DB.Debug().Where("username = ?", username).First(&userModel)
	if userModel.ID == 0 {
		return userDto, errors.New("未获取到当前用户信息")
	}
	userDto.Id = userModel.ID
	userDto.Name = userModel.Name
	userDto.Username = userModel.Username
	userDto.UpdateDate = userModel.UpdatedAt.Format("2006-01-02 15:04:05")
	userDto.CreateDate = userModel.CreatedAt.Format("2006-01-02 15:04:05")
	userDto.Avatar = userModel.Avatar
	userDto.State = userModel.State
	var userPrivilegeDtos []user.UserPrivilegeDto
	// 存入用户权限
	rows, _ := G.DB.Debug().Table("user_privilege").Select("(CASE privilege.type WHEN 0 THEN 'user' ELSE 'database' END) "+
		"AS privilegeType,privilege.code AS privilegeCode").
		Joins("left join privilege on privilege.id = user_privilege.privilege_id").Where("user_privilege.user_id = ?", userModel.ID).Rows()
	for rows.Next() {
		var userPrivilegeDto user.UserPrivilegeDto
		rows.Scan(&userPrivilegeDto.PrivilegeType, &userPrivilegeDto.PrivilegeCode)
		userPrivilegeDtos = append(userPrivilegeDtos, userPrivilegeDto)
	}
	userDto.UserPrivileges = userPrivilegeDtos
	return userDto, nil
}

// CheckUsernameAndPassword 校验用户名和密码，password为rsa加密后的密码
func (obj *UserService) CheckUsernameAndPassword(username string, password string) error {
	rsa := utils.NewRsa(G.C.Jwt.Rsa.PrivateKey)
	redisUserKey := fmt.Sprintf("%s_login_error", username)
	decryptPassword, err := rsa.Decrypt([]byte(password))
	if err != nil {
		// 此处输入错误的密码加密，也会记入密码错误次数中
		G.Logger.Errorf("解析密码失败，失败原因:[%s]，加密密码:[%s]", err, password)
		err = obj.checkLoginErrRedis(username, redisUserKey)
		if err != nil {
			return err
		}
	}
	// 校验现有的密码和argon2生成的加密密码正确性
	var user user2.User
	G.DB.Debug().Where("username = ?", username).First(&user)
	if user.ID == 0 {
		G.Logger.Errorf("用户名[%s]不存在", username)
		return errors.New("用户名不存在")
	}
	// 如果用户state为1，则用户已禁用
	if user.State == 1 {
		return errors.New("账户已锁定，请联系管理员")
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
	// 如果用户名密码匹配，则需要从redis中查看是否有密码输入错误的key，有的话则进行删除
	if matchPassword {
		getRedisUserKey, _ := G.RedisClient.Get(redisUserKey).Result()
		if getRedisUserKey != "" {
			G.RedisClient.Del(redisUserKey)
		}
		return nil
	}
	err = obj.checkLoginErrRedis(username, redisUserKey)
	if err != nil {
		return err
	}
	return nil
}

func (obj *UserService) checkLoginErrRedis(username string, redisKey string) error {
	// 需要先从redis获取，剩余输入次数，如果用户信息不存在，则将错误的用户名信息加入redis中，并提示剩余错误输入次数，次数到达将锁定账户，
	// 存在则进行次数递减,key格式是[用户名_login_error],值是5
	userLoginRedis, err := G.RedisClient.Get(redisKey).Result()
	if err != nil {
		if err != redis.Nil {
			G.Logger.Errorf("从redis中获取用户[%s]错误登录信息失败,redis key:[%s]", username, redisKey)
			return errors.New("系统错误")
		}
	}
	if userLoginRedis == "" {
		err = G.RedisClient.Set(redisKey, 5, 0).Err()
		if err != nil {
			G.Logger.Errorf("从redis中设置用户[%s]错误登录息信失败:[%s]", username, err.Error())
			return errors.New("系统错误")
		}
		return errors.New(fmt.Sprintf("用户名或密码不正确，剩余密码输入错误次数5次"))
	}
	if userLoginRedis == "1" {
		// 剩余一次机会时，进行账户锁定，并从redis中删除key
		affected := G.DB.Model(&user2.User{}).Where("username = ?", username).Update("state", 1).RowsAffected
		if affected == 0 {
			G.Logger.Errorf("禁用用户[%s]失败:[%s]", username, err.Error())
			return errors.New("系统错误")
		}
		err = G.RedisClient.Del(redisKey).Err()
		if err != nil {
			G.Logger.Errorf("从redis中删除用户[%s]key失败:[%s]", username, err.Error())
			return errors.New("系统错误")
		}
		return errors.New(fmt.Sprintf("账户已锁定，请联系管理员"))
	}
	if userLoginRedis > "1" && userLoginRedis <= "5" {
		err = G.RedisClient.Decr(redisKey).Err()
		if err != nil {
			G.Logger.Errorf("从redis中减少用户[%s]错误登录息信失败:[%s]", username, err.Error())
			return errors.New("系统错误")
		}
		getRedisUserKey, _ := G.RedisClient.Get(redisKey).Result()
		return errors.New(fmt.Sprintf("用户名或密码不正确，剩余密码输入错误次数%s次", getRedisUserKey))
	}
	return errors.New(fmt.Sprintf("用户名或密码不正确，剩余密码输入错误次数%s次", userLoginRedis))
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
