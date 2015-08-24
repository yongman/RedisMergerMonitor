package server

import (
	c "../conf"
	"github.com/gin-gonic/gin"
)

func RunHttpServer(meta *c.MonitorConf) {
	router := gin.Default()
	router.Static("/ui", "./public")
	router.Run(meta.HttpListen)
}
