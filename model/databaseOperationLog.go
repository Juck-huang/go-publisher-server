package model

import "gorm.io/gorm"

// DatabaseOperationLog 数据库操作日志表
type DatabaseOperationLog struct {
	gorm.Model
	Name       string `gorm:"type:varchar(155);comment:操作名称"`                   // 操作名称
	Type       int64  `gorm:"type:integer(2);comment:操作类型，0.新增，1.删除，2.修改，3.查询"` // 操作类型，（暂不记录查询）
	OperatorId int64  `gorm:"type:integer(100);comment:操作人id"`                  // 操作人id
	Detail     string `gorm:"type:varchar(2000);comment:详细内容"`                  // 详细内容,包含操作的sql
}

// TableName 自定义表名
func (DatabaseOperationLog) TableName() string {
	return "database_operation_log"
}
