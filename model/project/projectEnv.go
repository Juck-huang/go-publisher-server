package project

import "gorm.io/gorm"

type ProjectEnv struct {
	gorm.Model
	Name      string `gorm:"type:varchar(155);comment:项目环境名称"`
	ProjectId uint   `gorm:"type:integer(50);comment:项目id"`
}

// TableName 自定义表名
func (ProjectEnv) TableName() string {
	return "project_env"
}
