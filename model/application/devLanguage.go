package application

import "gorm.io/gorm"

// 开发语言
type DevLanguage struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);comment:名称"`   // 语言名称
	Code string `gorm:"type:varchar(100);comment:日志路径"` // 语言编码
}

// TableName 自定义表名
func (DevLanguage) TableName() string {
	return "dev_language"
}
