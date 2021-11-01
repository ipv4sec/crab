package middleware

import (
	"context"
	"crab/cluster"
	"crab/utils"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func DomainMiddleware(c *gin.Context) {
	if Memory.Domain != "" {
		c.Next()
	}

	domain, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(), "island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, "获取根域失败"))
		return
	}
	if domain.Data["root-domain"] == "example.com" {
		c.JSON(200, utils.SuccessResponse(map[string]interface{}{
			"todo": 1,
			"message": "未设置根域, 跳转到设置页面",
		}))
		return
	}
	c.Next()
}