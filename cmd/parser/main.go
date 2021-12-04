package main

import (
	"crab/parser"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"math/rand"
	"time"
)

func main() {
	var err error
	r := gin.Default()
	rand.Seed(time.Now().UnixNano())
	r.POST("/", parser.PostManifestHandlerFunc)

	err = r.Run(":4000")
	if err != nil {
		klog.Errorln("端口已被占用")
		panic(err)
	}
}
