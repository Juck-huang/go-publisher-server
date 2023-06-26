package model

import "gorm.io/gorm"

// ProjectType 项目类型
type ProjectType struct {
	gorm.Model
	Name      string `gorm:"type:varchar(155);comment:类型名称"`
	ProjectId int64  `gorm:"type:integer(50);comment:项目id"`
	ParentId  int64  `gorm:"type:integer(50);comment:父id"`
	IsLeaf    int64  `gorm:"type:boolean;comment:是否是叶子节点0.否1.是"`
	TreeId    string `gorm:"type:varchar(255);comment:树id"`
	TreeLevel int64  `gorm:"type:integer(20);comment:树级别"`
}

// TableName 自定义表名
func (ProjectType) TableName() string {
	return "project_type"
}
