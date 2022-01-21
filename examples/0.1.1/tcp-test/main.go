package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/mysql"
	"io/ioutil"
	"k8s.io/klog/v2"
	"time"
	"gopkg.in/yaml.v3"

)
type Mysql struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}
type Redis struct {
	Host     string
	Port     int
	Password string
}
type Config struct {
	Mysql Mysql
	Redis Redis
}
var db *gorm.DB

func main() {
	var err error
	var conf string
	flag.StringVar(&conf,"config", "config.yaml", "配置文件")
	flag.Parse()
	fmt.Println("读取配置文件")
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		panic(err)
	}
	klog.Infoln("解析配置文件")
	var cfg Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}
	route := gin.Default()


	route.GET("/redis", func(c *gin.Context) {
		// 建立连接
		conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port))
		if err != nil {
			fmt.Println("redis.Dial err=", err)
			return
		}
		// 通过go向redis写入数据 string [key - value]
		_, err = conn.Do("Set", "time", time.Now())
		if err != nil {
			fmt.Println("set err=", err)
			return
		}
		// 关闭连接
		defer conn.Close()
		// 读取数据 获取名字
		r, err := redis.String(conn.Do("Get", "time"))
		if err != nil {
			fmt.Println("set err=", err)
			return
		}
		c.JSON(200, gin.H{
			"from redis name: ": r,
		})
	})

	route.GET("/mysql", func(c *gin.Context) {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Mysql.Username,
			cfg.Mysql.Password,
			cfg.Mysql.Host,
			cfg.Mysql.Port,
			cfg.Mysql.Database)
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			fmt.Println("链接数据库失败")
		}
		klog.Infoln("连接到数据库")
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
	route.GET("/", func(c *gin.Context) {
		type r struct {
			Path string
			Method string
		}
		l := make([]r, 0)
		routers := route.Routes()
		for _, v := range routers {
			l = append(l, r{
				Path:    v.Path,
				Method: v.Method,
			})
		}
		c.JSON(200, gin.H{
			"api": l,
		})
	})
	err = route.Run(":3000")
	if err != nil {
		panic(err)
	}
}