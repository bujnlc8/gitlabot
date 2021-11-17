package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Getenv("BotDebug")) == 0 {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.POST("/", TransmitRobot)
	listenAddr := os.Getenv("listenAddr")
	if len(listenAddr) == 0 {
		listenAddr = "0.0.0.0:9090"
	}
	r.Run(listenAddr)
}
