package domain

import (
	"context"
	"crab/cluster"
	"crab/utils"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
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
		Result string `json:"result"`
	}{
		Status: 0,
		Result: "",
	}
	conf, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(),"island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域的键值对失败", err.Error())
		status.Status = 1
		status.Result = "保存根域失败"
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, status))
		return
	}
	conf.Data["root-domain"] = param.Value
	_, err = cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Update(context.Background(),conf, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorln("设置根域的键值对失败", err.Error())
		status.Status = 1
		status.Result = "保存根域失败"
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterSetConfigMap, status))
		return
	}
	status.Status = 2
	status.Result = "保存根域成功"
	c.JSON(200, utils.SuccessResponse(status))
}
