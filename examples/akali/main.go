package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"headers": c.Request.Header,
			"host": c.Request.Host,
		})
	})
	log.Fatal(r.Run(":3000"))
}