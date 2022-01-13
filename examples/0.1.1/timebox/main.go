package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"timestamp": time.Now().Unix(),
		})
	})
	err := route.Run(":3000")
	if err != nil {
		panic(err)
	}
}