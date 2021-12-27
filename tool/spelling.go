package tool

import (
	"crab/provider"
	"crab/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"k8s.io/klog/v2"
	"time"
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
	if param.Value == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	result, err := provider.Template()
	if err != nil {
		klog.Errorln("请求模板错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "请求模板错误"))
		return
	}
	timeNow := time.Now().Unix()
	saved := fmt.Sprintf("/tmp/%v.cue", timeNow)
	err = ioutil.WriteFile(saved, []byte(fmt.Sprintf("%s\n%s", result, param.Value)),0777)
	if err != nil {
		klog.Errorln("保存文件错误", saved, err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "保存文件错误"))
		return
	}
	cmd := fmt.Sprintf("cue vet %s", saved)
	output, _ := executor.ExecuteCommandWithCombinedOutput("bash", "-c", cmd)
	if output == "" {
		c.JSON(200, utils.SuccessResponse("正确"))
		return
	}
	c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, output))
}