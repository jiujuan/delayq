package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jiujuan/delayq/delayq"
)

func Route(router *gin.Engine) {
	router.POST("/add", delayq.AddJob)
	router.GET("/pop/", delayq.PopJob)
	router.GET("/delete/:id", delayq.DeleteJob)
}
