package user

import (
	"gorm.io/gorm"
)

// UserProject 用户项目表
type UserProject struct {
	gorm.Model
	UserId    uint `gorm:"type:varchar(50);comment:用户id"`
	ProjectId uint `gorm:"type:varchar(50);comment:项目id"`
}

// TableName 自定义表名
func (UserProject) TableName() string {
	return "user_project"
}
