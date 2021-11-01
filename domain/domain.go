package domain

import (
	"context"
	"crab/cluster"
	"crab/middleware"
	"crab/provider"
	"crab/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"strings"
)

func GetDomainHandlerFunc(c *gin.Context)  {
	domain, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(), "island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, "获取根域失败"))
		return
	}
	c.JSON(200, utils.SuccessResponse(domain.Data["root-domain"]))
}

func PutDomainHandlerFunc(c *gin.Context) {
	var param struct {
		Value string `json:"value"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil || param.Value == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	status := struct {
		Status int `json:"status"`
		Message string `json:"message"`
	}{
		Status: 0,
		Message: "未检测到域名的解析",
	}
	res, err := provider.HTTPClient.Get(fmt.Sprintf("http://%s/", param.Value), nil)
	if err != nil {
		klog.Errorln("请求该域名错误", err.Error())
		c.JSON(200, utils.SuccessResponse(status))
		return
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.Errorln("请求该域名错误", err.Error())
		c.JSON(200, utils.SuccessResponse(status))
		return
	}

	if  strings.TrimSpace(string(bytes)) != "crab" {
		klog.Errorln("接口返回域名错误", strings.TrimSpace(string(bytes)), param.Value)
		c.JSON(200, utils.SuccessResponse(status))
		return
	}
	conf, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(),"island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域的键值对失败", err.Error())
		status.Status = 1
		status.Message = "保存根域失败"
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, status))
		return
	}
	conf.Data["root-domain"] = param.Value
	_, err = cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Update(context.Background(),conf, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorln("设置根域的键值对失败", err.Error())
		status.Status = 1
		status.Message = "保存根域失败"
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterSetConfigMap, status))
		return
	}
	middleware.Memory.Domain = param.Value

	status.Status = 2
	status.Message = "保存根域成功"
	c.JSON(200, utils.SuccessResponse(status))
}
