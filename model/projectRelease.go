package model

import "gorm.io/gorm"

type ProjectRelease struct {
	gorm.Model
	Name            string `gorm:"type:varchar(155);comment:发布名称"`
	Description     string `gorm:"type:varchar(155);comment:描述"`
	ProjectId       int64  `gorm:"type:varchar(155);comment:项目id"`
	ProjectEnvId    int64  `gorm:"type:integer(50);comment:项目环境id"`
	ProjectTypeId   int64  `gorm:"type:integer(50);comment:项目类型id"`
	BuildScriptPath string `gorm:"type:varchar(155);comment:构建脚本路径"`
	Params          string `gorm:"type:varchar(155);comment:多参数，逗号分隔"` // 多参数，逗号分隔
}

// TableName 自定义表名
func (ProjectRelease) TableName() string {
	return "project_release"
}
