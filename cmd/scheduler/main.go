package main

import (
	"crab/cache"
	"crab/config"
	"crab/scheduler"
	"flag"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func main()  {
	var err error
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
	err = cache.Init(&cfg.Redis)
	if err != nil {
		panic(err)
	}

	go scheduler.Consumer()
	
	r := gin.Default()
	r.POST("/", scheduler.PostDeploymentHandlerFunc)

	err = r.Run(":3000")
	if err != nil {
		klog.Errorln("端口已被占用")
		panic(err)
	}
}