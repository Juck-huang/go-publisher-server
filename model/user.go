package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User 用户结构体
type User struct {
	gorm.Model
	Username         string    `gorm:"type:varchar(155);comment:用户名"` // 用户名
	Password         string    `gorm:"type:varchar(255)"`             // 密码，保存加密后的密码
	State            int       `gorm:"type:int(2)"`                   // 用户状态 0为正常， 1为禁用， 默认为0
	LastLoginDate    time.Time `gorm:"type:datetime"`                 // 上次登录日期
	LastLoginIp      string    `gorm:"type:varchar(255)"`             // 用户上次登录的ip地址
	LastLoginIpPlace string    `gorm:"type:varchar(255)"`             // 上次登录ip归属地
	Avatar           string    `gorm:"type:varchar(255)"`             // 个人头像地址
}
