package model

import "gorm.io/gorm"

type ProjectRelease struct {
	gorm.Model
	Name        string `gorm:"type:varchar(155);comment:发布项目名称"`
	Content     string `gorm:"type:varchar(155);comment:发布内容"`
	ProjectType int64  `gorm:"type:integer(50);comment:项目大类型"`
}

// TableName 自定义表名
func (ProjectRelease) TableName() string {
	return "project_release"
}
