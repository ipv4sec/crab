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
type item struct {
	Id int `json:"id"`
	Url string `json:"url"`
	Desc string `json:"desc"`
}
func main() {
	bytes, err := ioutil.ReadFile("/etc/configs/sso-alpha")
	if err != nil {
		klog.Errorln("读取依赖地址错误:", err.Error())
	}
	klog.Infoln("读取到的地址为:", string(bytes))
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		l := make([]item, 0)
		l = append(l, item{
			Id:   1,
			Url:  "/user",
			Desc: "从应用sso-alpha获取数据",
		})
		c.JSON(200, l)
	})
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