package user

import (
	"crab/cluster"
	"crab/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUserHandlerFunc(c *gin.Context)  {
	username := c.Param("username")
	if username == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}

	administrators, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(), "island-administrator", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取administrator的键值对失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, "获取用户信息失败"))
		return
	}
	v, ok := administrators.Data[username]
	if ok {
		c.JSON(200, utils.SuccessResponse(User{
			Username: username,
			Password: v,
		}))
		return
	}
	c.JSON(200, utils.SuccessResponse(struct {
		Error string `json:"error"`
	}{
		Error: "该用户不存在",
	}))
}
