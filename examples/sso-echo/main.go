package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gojek/heimdall/v7/httpclient"
	"io/ioutil"
	"k8s.io/klog/v2"
	"time"
)

var (
	HTTPClient = httpclient.NewClient(httpclient.WithHTTPTimeout(time.Second * 30))
)

func main() {
	ssoUrl, err := ioutil.ReadFile("/etc/configs/sso")
	if err != nil {
		klog.Errorln("读取依赖地址错误:", err.Error())
	}
	configs, err := ioutil.ReadFile("/etc/configs/userconfigs")
	if err != nil {
		klog.Errorln("读取运行时配置错误:", err.Error())
	}
	klog.Infoln("读取到的地址为:", string(ssoUrl))
	klog.Infoln("读取到的运行时配置为:", string(ssoUrl))
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		res, err := HTTPClient.Get(fmt.Sprintf("http://%s/user.json", string(ssoUrl)), nil)
		if err != nil {
			klog.Errorln("发送请求错误:", err.Error())
		}
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			klog.Errorln("读取返回数据错误:", err.Error())
		}
		c.String(200, fmt.Sprintf("运行时配置: %s, 从SSO读取到的数据: %s", string(configs), string(bodyBytes)))
	})
	err = route.Run(":80")
	if err != nil {
		panic(err)
	}
}