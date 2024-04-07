package main

import (
	"fmt"
	"net/http"
	"sprint/go/common/logger"
	"sprint/go/config"
	"sprint/go/data"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, data.CreateCommonSuccessResponse(gin.H{
			"message": "pong",
		}))
	})

	host, port := config.GetPortAndHost()
	serverUrl := fmt.Sprintf("%s:%s", host, port)

	logger.Info("Server started on ", serverUrl)
	r.Run(serverUrl) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
