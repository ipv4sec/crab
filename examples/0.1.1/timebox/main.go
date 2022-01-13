package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	var Memory = []int64{time.Now().Unix()}

	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})
	route.GET("/timestamp.json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"timestamp": Memory,
		})
	})
	route.POST("/timestamp.json", func(c *gin.Context) {
		Memory = append(Memory, time.Now().Unix())
		c.JSON(200, gin.H{
			"timestamp": Memory,
		})
	})
	err := route.Run(":3000")
	if err != nil {
		panic(err)
	}
}