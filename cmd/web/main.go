package main

import (
	"crab/app"
	"crab/cluster"
	"crab/config"
	"crab/db"
	"crab/status"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func main() {
	var err error
	var conf string
	flag.StringVar(&conf, "config", "config.yaml", "配置文件")
	flag.Parse()

	klog.Infoln("读取配置文件")
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		panic(err)
	}

	klog.Infoln("解析配置文件")
	var cfg config.Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}

	klog.Infoln("连接到数据库")
	err = db.Init(&cfg.Mysql)
	if err != nil {
		panic(err)
	}

	if false {
		err = db.Client.AutoMigrate(&status.Status{})
		if err != nil {
			panic(err)
		}
		err = db.Client.AutoMigrate(&app.App{})
		if err != nil {
			panic(err)
		}
	}

	klog.Infoln("开始集群认证")
	err = cluster.Init()
	if err != nil {
		panic(fmt.Errorf("获取集群认证失败: %w", err))
	}
	klog.Infoln("集群认证成功")

	klog.Infoln("开始提供服务")
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	routers := r.Group("/api")
	{
		routers.GET("/app", app.GetAppHandlerFunc)
		routers.PUT("/app", app.PutAppHandlerFunc)
		routers.POST("/app", app.PostAppHandlerFunc)
		routers.DELETE("/app", app.DeleteAppHandlerFunc)
	}
	err = r.Run("0.0.0.0:3000")
	if err != nil {
		panic(fmt.Errorf("监听端口失败: %w", err))
	}
}