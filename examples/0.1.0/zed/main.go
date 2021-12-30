package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gojek/heimdall/v7/httpclient"
	"io/ioutil"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

var (
	HTTPClient = httpclient.NewClient(httpclient.WithHTTPTimeout(time.Second * 30))
)

func main() {
	route := gin.Default()
	user := route.Group("/")
	{
		user.GET("/", func(c *gin.Context) {
			files, err := ioutil.ReadDir("/opt/")
			if err != nil {
				klog.Errorln("读取依赖地址错误:", err.Error())
				return
			}
			vals := []string{}
			for i := 0; i < len(files); i++ {
				bytes, err := ioutil.ReadFile("/opt/"+ files[i].Name())
				if err != nil {
					vals = append(vals, err.Error())
				} else {
					vals = append(vals, string(bytes))
				}
			}
			c.String(200, strings.Join(vals, ","))
		})
		user.GET("/user", func(c *gin.Context) {
			res, err := HTTPClient.Get(fmt.Sprintf("http://%s", c.Query("a")), nil)
			if err != nil {
				klog.Errorln("发送请求错误:", err.Error())
			}
			bodyBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				klog.Errorln("读取返回数据错误:", err.Error())
			}
			c.String(200, fmt.Sprintf("依赖: %s", string(bodyBytes)))
		})
	}
	err := route.Run(":80")
	if err != nil {
		panic(err)
	}
}