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
	value, ok := administrators.Data[username]
	if ok {
		c.JSON(200, utils.RowResponse(User{
			Username: username,
			Password: value,
		}))
		return
	}
	c.JSON(200, utils.RowResponse(map[string]interface{}{"error":"该用户不存在"}))
}

func PutUserHandlerFunc(c *gin.Context) {
	username := c.Param("username")
	if username != "root" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	var param struct{
		Password string `json:"password"`
		OldPassword string `json:"oldPassword"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		klog.Errorln("参数错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	cm, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(), "island-administrator", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取用户信息失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "获取用户信息失败"))
		return
	}
	if param.OldPassword != cm.Data["root"] {
		klog.Errorln("密码不符", param.OldPassword, cm.Data["root"])
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "密码不符"))
		return
	}
	cm.Data["root"] = param.Password
	_, err = cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Update(context.Background(), cm, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorln("更新用户信息失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterUpdateConfigMap, "更新用户信息失败"))
		return
	}
	c.JSON(200, utils.SuccessResponse("设置成功"))
}
