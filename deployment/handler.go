package deployment

import (
	"crab/utils"
	"github.com/gin-gonic/gin"
)

func PutDeploymentHandlerFunc(c *gin.Context) {
	var param struct {
		Deployment string
		Parameters string
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse("部署成功"))
}