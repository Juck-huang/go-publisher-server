package service

import (
	"errors"
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
func (obj UserService) CheckUsernameAndPassword(username string, password string) error {
	rsa := utils.NewRsa(G.C.Rsa.PrivateKey)
	decrypt, err := rsa.Decrypt([]byte(password))
	if err != nil {
		G.Logger.Errorf("登录失败，失败原因:[%s]", err)
		return errors.New("登录失败，请联系管理员")
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
