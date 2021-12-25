package tool

import (
	"crab/provider"
	"crab/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func GetTemplateHandlerFunc(c *gin.Context) {
	result, err := provider.Template()
	if err != nil {
		klog.Errorln("请求模板错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "请求模板错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(result))
}
