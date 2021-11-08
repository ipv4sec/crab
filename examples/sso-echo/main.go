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

	bytes, err := ioutil.ReadFile("/etc/configs/sso-alpha.host")
	if err != nil {
		klog.Errorln("读取依赖地址错误:", err.Error())
	}
	klog.Infoln("读取到的地址为:", string(bytes))
	route := gin.Default()
	user := route.Group("/user")
	{
		user.GET("/", func(c *gin.Context) {
			res, err := HTTPClient.Get(fmt.Sprintf("http://%s/user.json", string(bytes)), nil)
			if err != nil {
				klog.Errorln("发送请求错误:", err.Error())
			}
			bodyBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				klog.Errorln("读取返回数据错误:", err.Error())
			}
			c.String(200, fmt.Sprintf("sso: %s", string(bodyBytes)))
		})
	}
	err = route.Run(":80")
	if err != nil {
		panic(err)
	}
}