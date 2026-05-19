package main

import (
	"github.com/gilangages/kopi-popi/config"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config.ConnectDB()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "api running",
		})
	})

	r.Run(":8080")
}
