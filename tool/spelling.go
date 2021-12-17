package tool

import (
	"crab/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func PostSpellingHandlerFunc(c *gin.Context) {
	var param struct {
		Value string `json:"value"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		klog.Errorln("参数错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	// TODO
	c.JSON(200, utils.SuccessResponse(param.Value))
}