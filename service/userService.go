package service

import (
	"errors"
	"hy.juck.com/go-publisher-server/config"
)

var (
	G = config.G
)

type UserService struct {
}

// CheckUsernameAndPassword 校验用户名和密码
func (obj UserService) CheckUsernameAndPassword(username string, password string) error {
	prepare, err := G.DB.Prepare("select count(1) from user where username = ? and password = ?")
	if err != nil {
		return err
	}
	var num int64
	err = prepare.QueryRow(username, password).Scan(&num)
	if err != nil {
		return err
	}
	if num > 0 {
		return nil
	}
	return errors.New("用户名或密码不正确")
}
