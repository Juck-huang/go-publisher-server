package server

import "gorm.io/gorm"

// 服务器信息
type ServerInfo struct {
	gorm.Model
	Name      string `gorm:"type:varchar(155);comment:服务器名称"`
	State     int    `gorm:"type:integer(2);default:0;comment:服务器是否启用,0.未启用,1.启用"`
	ProjectId string `gorm:"type:varchar(50);comment:所属项目id"`
	UId       string `gorm:"type:varchar(255);comment:服务器唯一标识"`
}

// TableName 自定义表名
func (ServerInfo) TableName() string {
	return "server_info"
}
