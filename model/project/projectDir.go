package project

import "gorm.io/gorm"

type ProjectDir struct {
	gorm.Model
	Name          string `gorm:"type:varchar(155);comment:项目目录名称"`
	ProjectId     uint   `gorm:"type:integer(50);comment:项目id"`
	ProjectEnvId  int64  `gorm:"type:integer(50);comment:项目环境id"`
	ProjectTypeId int64  `gorm:"type:integer(50);comment:项目类型id"`
	ProjectPath   string `gorm:"type:varchar(155);comment:项目目录根路径"`
}

// TableName 自定义表名
func (ProjectDir) TableName() string {
	return "project_dir"
}
