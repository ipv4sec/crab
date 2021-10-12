package main

import (
	"crab/config"
	"crab/db"
	"crab/manifest"
	"flag"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func Init(){
	var conf string
	flag.StringVar(&conf, "config", "config.yaml", "配置文件")
	flag.Parse()
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		panic(err)
	}

	var cfg config.Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}

	err = db.Init(&cfg.Mysql)
	if err != nil {
		panic(err)
	}
}

func main(){
	var err error
	Init()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//解析依赖
	r.POST("/", manifest.PostManifestHandlerFunc)
	//生成k8s文件
	r.PUT("/", manifest.PutManifestHandlerFunc)

	err = r.Run(":3000")
	if err != nil {
		klog.Errorln("端口已被占用")
		panic(err)
	}
}
