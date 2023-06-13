package model

// DatabaseOperationLog 数据库操作日志表
type DatabaseOperationLog struct {
	Id               int    `json:"id"`               // id
	Name             string `json:"name"`             // 操作名称
	RecordCreateDate string `json:"recordCreateDate"` // 记录创建时间
	RecordUpdateDate string `json:"recordUpdateDate"` // 记录更新时间
	IpAddress        string `json:"ipAddress"`        // 操作人的ip地址
}
