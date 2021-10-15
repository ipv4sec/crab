package main

import (
	"crab/manifest"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func main() {
	var err error
	r := gin.Default()

	//parse manifest.yaml to k8s.yaml
	r.POST("/", manifest.PostManifestHandlerFunc)

	err = r.Run(":3000")
	if err != nil {
		klog.Errorln("端口已被占用")
		panic(err)
	}
}
