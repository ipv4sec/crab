package mirror

import (
	"context"
	"crab/cluster"
	"crab/utils"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetMirrorHandlerFunc(c *gin.Context)  {
	domain, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(), "island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, "获取根域失败"))
		return
	}
	c.JSON(200, utils.SuccessResponse(domain.Data["mirror"]))
}

func PutMirrorHandlerFunc(c *gin.Context) {
	var param struct {
		Value string `json:"mirror"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil || param.Value == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	conf, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(),"island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域的键值对失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, "获取根域的键值对失败"))
		return
	}
	conf.Data["mirror"] = param.Value
	_, err = cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Update(context.Background(),conf, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorln("设置根域的键值对失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterSetConfigMap, "设置根域的键值对失败"))
		return
	}
	c.JSON(200, utils.SuccessResponse("设置成功"))
}
