package database

type TotalDto struct {
	DbName       string   `json:"dbName"`       // 数据库名称
	IgnoreTables []string `json:"ignoreTables"` // 需要忽略的表
	Type         int      `json:"type"`         // 1.是备份数据库，2.是导出数据库
}

// 存入redis的结构体
type RedisDataDto struct {
	Completed  bool   `json:"completed"`
	CreateTime string `json:"create_time"`
}
