package config

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() {
	__initSqlite()
	__initTable()
}

// 初始化本地sqlite数据库
func __initSqlite() {
	var err error
	G.DB, err = sql.Open("sqlite3", G.C.DB.Sqlite3.Path)
	if err != nil {
		panic(fmt.Sprintf("加载本地数据库失败，失败原因%s", err))
	}
	// 测试连接是否成功
	err = G.DB.Ping()
	if err != nil {
		panic(fmt.Sprintf("加载本地数据库失败，失败原因%s", err))
	}
	G.Logger.Info("加载本地sqlite数据库成功")
}

func __initTable() {
	_, err := G.DB.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			username TEXT,
			password TEXT,
			create_date TEXT,
		    last_login_date TEXT
		)
	`)
	if err != nil {
		panic(fmt.Sprintf("初始化本地用户表失败，失败原因%s", err.Error()))
	}
	G.Logger.Info("初始化本地表成功")
}
