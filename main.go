package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jiujuan/delayq/config"
	"github.com/jiujuan/delayq/logger"
	"github.com/jiujuan/delayq/router"
)

func main() {
	config.InitConfig()
	logger.InitLog(config.QConfig.Log)
	logger.Logger.Info("init config, logger ")

	app := gin.New()
	router.Route(app)
	// host := fmt.Sprintf("%s:%d", config.QConfig.App.IP, config.QConfig.App.Port)
	// logger.Logger.Info("read application config: " + host)

	host := fmt.Sprintf(":%d", config.QConfig.App.Port)
	logger.Logger.Info("init Application, host is " + host)
	app.Run(host)
}
