package database

type SimpleDto struct {
	DbName       string   `json:"dbName"`       // 数据库名称
	ExportTables []string `json:"exportTables"` // 需要导出的表
}
