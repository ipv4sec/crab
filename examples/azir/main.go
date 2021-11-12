package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	PK     int64  `gorm:"primaryKey" json:"-"`
	Name   string
}

func (User) TableName() string {
	return "ttt"
}

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "123456", "127.0.0.1", "3306", "ttt")
	Client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/:id", func(c *gin.Context) {
		var user User
		var err error
		err = Client.Model(&User{}).Where("pk = ?", c.Param("id")).Find(&user).Error
		if err != nil {
			panic(err)
		}
		var v interface{}
		err = json.Unmarshal([]byte(user.Name), &v)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{
			"name": v,
		})
	})
	r.POST("/", func(c *gin.Context) {
		var param struct{
			Name interface{} `json:"name"`
		}
		var err error
		err = c.ShouldBindJSON(&param)
		if err != nil {
			panic(err)
		}
		s, err := json.Marshal(param.Name)
		if err != nil {
			panic(err)
		}
		user := &User{
			Name: string(s),
		}
		err = Client.Create(user).Error
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{
			"pk": user.PK,
		})
	})
	r.Run(":3000")
}