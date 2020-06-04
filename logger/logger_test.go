package logger

import (
	"testing"
	"time"

	"github.com/jiujuan/delayq/config"

	"go.uber.org/zap"
)

func TestZapLog(t *testing.T) {
	config, _ := config.LoadConfig()
	t.Log(config)
	logger := InitZapConfig(config.Log)
	// defer logger.Sync()
	logger.Info("log 初始化", zap.String("url", "www.baidu.com"), zap.Duration("back", time.Second))
	logger.Error("log error", zap.String("error", "connect faile"))
}

func TestInitZapConfig(t *testing.T) {
	config, _ := config.LoadConfig()
	t.Log(config, config.Log)

	zaplog := InitZapConfig(config.Log)
	t.Log(zaplog)
}
