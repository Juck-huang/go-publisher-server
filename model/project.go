package model

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name string `gorm:"type:varchar(155);comment:项目名称"`
}

// TableName 自定义表名
func (Project) TableName() string {
	return "project"
}
