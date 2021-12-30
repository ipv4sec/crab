package tool

import (
	"bytes"
	"crab/provider"
	"crab/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os/exec"
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
	cmd := exec.Command("bash","-c", fmt.Sprintf("cue eval %s", saved))
	var erro bytes.Buffer
	cmd.Stderr = &erro
	err = cmd.Run()
	if err != nil {
		klog.Errorln("执行命令错误", err.Error(), string(erro.Bytes()))
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, string(erro.Bytes())))
		return
	}
	if string(erro.Bytes()) == "" {
		c.JSON(200, utils.SuccessResponse("正确"))
		return
	}
	c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, string(erro.Bytes())))
}