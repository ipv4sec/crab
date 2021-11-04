package main

import (
	"crab/app"
	"crab/cluster"
	"crab/config"
	"crab/db"
	"crab/deployment"
	"crab/domain"
	"crab/mirror"
	"crab/status"
	"crab/storage"
	"crab/user"
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
	flag.StringVar(&conf,"config", "config.yaml", "配置文件")
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

	klog.Infoln("开始集群认证")
	err = cluster.Init()
	if err != nil {
		panic(fmt.Errorf("获取集群认证失败: %w", err))
	}
	klog.Infoln("集群认证成功")

	klog.Infoln("开始提供服务")
	gin.SetMode(gin.ReleaseMode)
	routers := gin.Default()
	routers.GET("/", func(c *gin.Context) {
		c.String(200, "crab")
	})

	routers.GET("/user/:username", user.GetUserHandlerFunc)
	routers.PUT("/user/:username", user.PutUserHandlerFunc)

	routers.POST("/app", app.PostAppHandlerFunc)
	routers.GET("/app", app.GetAppsHandlerFunc)
	routers.GET("/app/:id", app.GetAppHandlerFunc)
	routers.GET("/app/:id/status", app.GetAppStatusHandlerFunc)
	routers.PUT("/app/:id", app.PutAppHandlerFunc)
	routers.DELETE("/app/:id", app.DeleteAppHandlerFunc)

	routers.GET("/cluster/addrs", storage.GetAddrsHandlerFunc)
	routers.GET("/cluster/domain", domain.GetDomainHandlerFunc)
	routers.PUT("/cluster/domain", domain.PutDomainHandlerFunc)

	routers.GET("/cluster/mirror", mirror.GetMirrorHandlerFunc)
	routers.PUT("/cluster/mirror", mirror.PutMirrorHandlerFunc)

	routers.PUT("/deployment", deployment.PutDeploymentHandlerFunc)

	routers.GET("/status/:id", status.GetStatusHandlerFunc)
	routers.GET("/status/:id/:componentName", status.GetComponentStatusHandlerFunc)
	routers.PUT("/status/:id/:componentName/:statusCode", status.PutStatusHandlerFunc)


	err = routers.Run("0.0.0.0:3000")
	if err != nil {
		panic(fmt.Errorf("监听端口失败: %w", err))
	}
}