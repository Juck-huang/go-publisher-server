package user

import (
	"gorm.io/gorm"
)

// UserProjectEnv 用户项目环境结构体
type UserProjectEnv struct {
	gorm.Model
	UserId       uint `gorm:"type:varchar(50);comment:用户id"`
	ProjectId    uint `gorm:"type:varchar(50);comment:项目id"`
	ProjectEnvId uint `gorm:"type:varchar(50);comment:项目环境id"`
}

// TableName 自定义表名
func (UserProjectEnv) TableName() string {
	return "user_project_env"
}
