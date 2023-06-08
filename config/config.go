package config

import (
	"database/sql"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 配置文件指针
var config = new(Config)

// G 全局配置调用
var G = new(Global)

type Global struct {
	C      *Config
	Logger *zap.SugaredLogger
	DB     *sql.DB
}

type Config struct {
	Server struct {
		Port int64 `mapstructure:"port"`
	} `mapstructure:"server"`
	DB struct {
		Sqlite3 struct {
			Path string `mapstructure:"path"`
		} `mapstructure:"sqlite3"`
		Mysql struct {
			Host       string `json:"host"`
			Port       string `json:"port"`
			Username   string `json:"username"`
			Password   string `json:"password"`
			BackUpPath string `json:"backUpPath"`
		}
	} `mapstructure:"db"`
	Zap struct {
		FileName   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"maxsize"`
		MaxBackups int    `mapstructure:"max-backups"`
		MaxAge     int    `mapstructure:"max-age"`
		Compress   bool   `mapstructure:"compress"`
		Mode       string `mapstructure:"mode"`
		Level      string `mapstructure:"level"`
	} `mapstructure:"zap"`
}

func InitLog() {
	logger := NewLogger()
	G.Logger = logger.sugarLogger
}

func InitViper() {
	// 指定配置文件路径
	viper.SetConfigFile("./config.yaml")
	// 读取配置文件信息
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err = viper.Unmarshal(config); err != nil {
		panic(fmt.Errorf("解析配置文件失败:%s", err))
	}
	G.C = config
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		G.Logger.Info("配置文件被修改了")
		if err = viper.Unmarshal(config); err != nil {
			panic(fmt.Errorf("解析配置文件失败:%s", err))
		}
	})
}

func init() {
	InitViper()
	InitLog()
	InitDB()
	G.Logger.Info("初始化所有配置成功")
}
