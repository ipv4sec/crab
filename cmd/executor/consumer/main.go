package main

import (
	"crab/cache"
	"crab/config"
	"crab/scheduler"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func main()  {
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

	go scheduler.Consumption()
	
	r := gin.Default()
	r.POST("/", scheduler.PostDeploymentHandlerFunc)

	err = r.Run(":3000")
	if err != nil {
		panic(fmt.Errorf("端口被占用: %w", err))
	}
}