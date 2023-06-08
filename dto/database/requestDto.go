package database

type RequestDto struct {
	DbName       string   `json:"dbName"`       // 数据库名称
	IgnoreTables []string `json:"ignoreTables"` // 需要忽略的表
	Status       int      `json:"status"`       // 1.是备份数据库，2.是导出数据库
}
