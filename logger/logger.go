package logger

import (
	"github.com/jiujuan/delayq/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLog 初始化日志
func InitLog(conf config.LogConfig) {
	Logger = InitZapConfig(conf)
}

// InitZapConfig 初始化zap配置
func InitZapConfig(conf config.LogConfig) *zap.Logger {
	zapcfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(getLevel(conf.LogLevel)),
		OutputPaths:      []string{conf.AccessLog},
		ErrorOutputPaths: []string{conf.ErrorLog},
		Encoding:         conf.LogEncode,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "time",
			CallerKey:     "line",
			StacktraceKey: "trace",

			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
			EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
			EncodeName:     zapcore.FullNameEncoder,
		},
	}

	logger, err := zapcfg.Build()
	if err != nil {
		panic("zap config init fail:" + err.Error())
	}

	return logger
}

func getLevel(level string) zapcore.Level {
	var l zapcore.Level
	switch level {
	case "debug":
		l = zapcore.DebugLevel
	case "info":
		l = zapcore.InfoLevel
	case "warn":
		l = zapcore.WarnLevel
	case "error":
		l = zapcore.ErrorLevel
	case "panic":
		l = zapcore.PanicLevel
	case "fatal":
		l = zapcore.FatalLevel
	default:
		l = zapcore.InfoLevel

	}
	return l
}
