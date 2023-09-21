package user

import "gorm.io/gorm"

// UserPrivilege 用户权限结构体
type UserPrivilege struct {
	gorm.Model
	UserId      int `gorm:"type:int(50);comment:用户id""`
	PrivilegeId int `gorm:"type:int(50);comment:权限id""`
}

// TableName 自定义表名
func (UserPrivilege) TableName() string {
	return "user_privilege"
}
