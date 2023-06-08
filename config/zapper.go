package config

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Logger struct {
	sugarLogger *zap.SugaredLogger
	mode        string
}

func NewLogger() *Logger {

	level := G.C.Zap.Level
	var levelEnable zapcore.Level
	switch level {
	case "error":
		levelEnable = zap.ErrorLevel
	case "info":
		levelEnable = zapcore.InfoLevel
	case "debug":
		levelEnable = zapcore.DebugLevel
	case "warn":
		levelEnable = zapcore.WarnLevel
	default:
	}

	var logger = &Logger{}
	logger.mode = G.C.Zap.Mode
	writeSyncer := logger.getLogWriter()
	encoder := logger.getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, levelEnable)
	caller := zap.New(core, zap.AddCaller())
	sugarLogger := caller.Sugar()
	logger.sugarLogger = sugarLogger

	return logger
}

func (o *Logger) getEncoder() zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	if o.mode == "dev" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (o *Logger) getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   G.C.Zap.FileName,   // 日志保存文件地址
		MaxSize:    G.C.Zap.MaxSize,    // 日志文件达到多少mb开始备份
		MaxBackups: G.C.Zap.MaxBackups, // 备份的文件数量
		MaxAge:     G.C.Zap.MaxAge,     // 最大备份保留天数
		Compress:   G.C.Zap.Compress,   // 是否压缩备份的日志
	}
	// 开发环境把日志文件输出到终端上,生产环境把日志记录到日志文件中
	if o.mode == "dev" {
		ws := io.MultiWriter(os.Stdout)
		return zapcore.AddSync(ws)
	} else {
		return zapcore.AddSync(lumberJackLogger)
	}
}
