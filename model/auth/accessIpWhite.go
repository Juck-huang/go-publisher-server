package auth

import "gorm.io/gorm"

// 访问ip白名单
type AccessIpWhite struct {
	gorm.Model
	ProjectId string `gorm:"type:varchar(100);comment:项目id"`      // 项目id
	ServerUId string `gorm:"type:varchar(200);comment:服务器唯一标识id"` // 服务器唯一标识id
	IpList    string `gorm:"type:varchar(255);comment:公网ip地址列表"`  // 服务器公网ip地址列表
}

// TableName 自定义表名
func (AccessIpWhite) TableName() string {
	return "access_ip_white"
}
