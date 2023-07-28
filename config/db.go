package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	databaseTable "hy.juck.com/go-publisher-server/model/database"
	"hy.juck.com/go-publisher-server/model/project"
	"hy.juck.com/go-publisher-server/model/user"
)

func InitDB() {
	initMysql()
}

// 初始化mysql
func initMysql() {

	mysqlConfig := G.C.Application.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Db)
	database, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 255,
	}), &gorm.Config{})

	if err != nil {
		panic("初始化数据库错误：" + err.Error())
	}

	sqlDB, err := database.DB()
	if err != nil {
		panic("初始化数据库错误：" + err.Error())
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	// sqlDB.SetConnMaxLifetime(time.Hour)
	G.DB = database
	// 自动迁移表
	database.AutoMigrate(
		&user.User{},
		&user.Privilege{},
		&user.UserPrivilege{},
		&project.Project{},
		&databaseTable.DatabaseOperationLog{},
		&project.ProjectRelease{},
		&project.ProjectType{},
		&project.ProjectEnv{},
		&project.ProjectDir{},
	)
	G.Logger.Infof("初始化数据库连接成功")
}
