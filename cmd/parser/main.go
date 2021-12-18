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

	//mainifest.yaml翻译为k8s资源文件
	r.POST("/", parser.PostManifestHandlerFunc)

	//获取k8s自动生成的cue模板
	r.GET("/systemTemplate", parser.GetSystemTemplateFunc)

	err = r.Run(":4000")
	if err != nil {
		klog.Errorln("端口已被占用")
		panic(err)
	}
}
