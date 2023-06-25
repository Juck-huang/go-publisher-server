package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户结构体
type User struct {
	gorm.Model
	Name             string    `gorm:"type:varchar(155);comment:用户名称"`
	Username         string    `gorm:"type:varchar(155);comment:用户名"`                        // 用户名
	Password         string    `gorm:"type:varchar(255)"`                                    // 密码，保存argon2加密后的密码
	State            int       `gorm:"type:int(2);default:0;comment:用户状态 0为正常， 1为禁用， 默认为0""` // 用户状态 0为正常， 1为禁用， 默认为0
	LastLoginDate    time.Time `gorm:"type:datetime(3);comment:上次登录日期"`                      // 上次登录日期
	LastLoginIp      string    `gorm:"type:varchar(255);comment:用户上次登录的ip地址"`                // 用户上次登录的ip地址
	LastLoginIpPlace string    `gorm:"type:varchar(255);comment:上次登录ip归属地"`                  // 上次登录ip归属地
	Avatar           string    `gorm:"type:varchar(255);comment:个人头像地址"`                     // 个人头像地址
}

// TableName 自定义表名
func (User) TableName() string {
	return "user"
}
