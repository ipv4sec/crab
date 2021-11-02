package status

import (
	"crab/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func GetStatusHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(1))
}

func GetComponentStatusHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	componentName := c.Param("componentName")
	if id == "" || componentName == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	klog.Info("status:", id, componentName)
	c.JSON(200, utils.SuccessResponse(1))
}

func PutStatusHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	componentName := c.Param("componentName")
	statusCode := c.Param("statusCode")
	if id == "" || componentName == "" || statusCode == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	var param struct {
		Message string
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}

	klog.Info("status:", id, componentName, statusCode, param.Message)
	c.JSON(200, utils.SuccessResponse("设置组件状态成功"))
}