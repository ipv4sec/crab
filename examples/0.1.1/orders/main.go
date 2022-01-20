package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func main() {
	db, err := gorm.Open(
		mysql.Open("root:123456@tcp(db:3306)/orders?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		var order struct {
			PK        uint `gorm:"primarykey"`
			ID        string
			Status    int
			CreateAt time.Time
			UpdateAt time.Time
		}
		db.Table("orders").First(&order)
		c.JSON(200, gin.H{
			"order": order,
		})
	})
	err = route.Run(":3000")
	if err != nil {
		panic(err)
	}
}