package application

import "gorm.io/gorm"

// 应用表
type Application struct {
	gorm.Model
	Name          string `gorm:"type:varchar(100);comment:名称"`     // 名称
	Code          string `gorm:"type:varchar(100);comment:编码"`     // 编码
	LogPath       string `gorm:"type:varchar(100);comment:日志路径"`   // 日志路径
	ConfigPath    string `gorm:"type:varchar(100);comment:配置文件路径"` // 配置文件路径，支持多个配置文件路径用逗号分隔
	AppPath       string `gorm:"type:varchar(100);comment:应用程序路径"` // 应用程序路径
	DevLanguageId uint   `gorm:"type:varchar(100);comment:开发语言id"` // 开发语言id
	ScriptPath    string `gorm:"type:varchar(100);comment:脚本路径"`   // 脚本路径
	Params        string `gorm:"type:varchar(100);comment:参数"`     // 参数
}

// TableName 自定义表名
func (Application) TableName() string {
	return "application"
}
