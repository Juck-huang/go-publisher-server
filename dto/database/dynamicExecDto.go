package database

type DynamicExecDto struct {
	DbName string `json:"dbName"` // 数据库名称
	Sql    string `json:"sql"`    // 需要执行的sql
}
