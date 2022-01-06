package deployment

import (
	"crab/aam/v1alpha1"
	"crab/app"
	"crab/db"
	"crab/exec"
	"crab/provider"
	"crab/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/klog/v2"
	"time"
)

var (
	executor = exec.CommandExecutor{}
)

func PutDeploymentHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	manifestFileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "接收文件错误"))
		return
	}
	manifestFile, err := manifestFileHeader.Open()
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "打开文件错误"))
		return
	}
	manifestBytes, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "读取文件错误"))
		return
	}

	var manifest v1alpha1.Application
	err = yaml.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		klog.Errorln("解析描述文件错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "解析描述文件错误"))
		return
	}

	var ins app.App
	err = db.Client.Where("id = ?", id).Find(&ins).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "该实例不存在"))
		return
	}

	var vals []struct{
		Name string `json:"name"`

		ID string `json:"id"`
		Location string `json:"location"`

		EntryService string
	}
	err = json.Unmarshal([]byte(ins.Additional), &vals)
	if err != nil {
		klog.Errorln("序列化依赖失败", err.Error())
	}
	var parameters interface{}
	err = json.Unmarshal([]byte(ins.Parameters), &parameters)
	if err != nil {
		klog.Errorln("序列化运行时配置失败", err.Error())
		parameters = ""
	}

	val, err1, err2 := provider.Yaml(string(manifestBytes), ins.ID, ins.Entry, parameters, provider.ConvertToDependency(vals))
	if err1 != nil {
		klog.Errorln("连接到翻译器错误:", err1.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, err1.Error()))
		return
	}
	if err2 != nil {
		klog.Errorln("连接到翻译器错误:", err2.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "连接到翻译器错误"))
		return
	}
	klog.Infoln("要执行的文件内容为:", val)
	timeNow := time.Now().Unix()
	saved := fmt.Sprintf("/tmp/%v.yaml", timeNow)
	err = ioutil.WriteFile(saved, []byte(val),0777)
	if err != nil {
		klog.Errorln("保存文件错误", saved, err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "保存文件错误"))
	}
	command := fmt.Sprintf("/usr/local/bin/kubectl apply -f %s", saved)
	output, _ := executor.ExecuteCommandWithCombinedOutput("bash", "-c", command)
	klog.Infoln("执行命令结果:", output)
	c.JSON(200, utils.SuccessResponse("部署成功"))
}