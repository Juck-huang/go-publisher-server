package user

import (
	"gorm.io/gorm"
)

// Privilege 权限结构体
type Privilege struct {
	gorm.Model
	Name string `gorm:"type:varchar(155);comment:权限名称"`                 // 权限名称
	Code string `gorm:"type:varchar(155);comment:权限编码"`                 // 权限编码
	Type int    `gorm:"type:int(2);default:0;comment:所属类型:0.用户，1.数据库;"` // 所属类型，0.用户，1.数据库
}

// TableName 自定义表名
func (Privilege) TableName() string {
	return "privilege"
}
