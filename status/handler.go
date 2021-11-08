package status

import (
	"crab/aam/v1alpha1"
	"crab/app"
	"crab/db"
	"crab/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
	"strconv"
)

func GetComponentStatusHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	componentName := c.Param("componentName")
	if id == "" || componentName == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var status Status
	err := db.Client.Where("id = ? AND name = ?", id, componentName).Find(&status).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库查询错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(status.Status))
}

func PutStatusHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	componentName := c.Param("componentName")
	statusCode, err := strconv.Atoi(c.Param("statusCode"))
	if id == "" || componentName == "" || err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var param struct {
		Message string
	}
	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	status := Status{
		ID:        id,
		Name: componentName,
		Status:    statusCode,
		Message:   param.Message,
	}
	err = db.Client.Create(&status).Error
	if err != nil {
		klog.Errorln("数据库保存错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库保存错误"))
		return
	}
	var val app.App
	err = db.Client.Where("id = ?", id).Find(&val).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseBadRequest, "该实例不存在"))
		return
	}
	var manifest v1alpha1.Application
	err = yaml.Unmarshal([]byte(val.Manifest), &manifest)
	if err != nil {
		klog.Errorln("解析描述文件错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "解析描述文件错误"))
		return
	}
	var componentStatus []Status
	err = db.Client.Where("id = ?", id).Find(&componentStatus).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库查询错误"))
		return
	}
	if len(manifest.Spec.Workloads) == len(componentStatus) {
		err = db.Client.Model(app.App{}).Where("id = ?", id).Update("status", 2).Error
		if err != nil {
			klog.Errorln("数据库更新错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "更新状态错误"))
			return
		}
	}
	c.JSON(200, utils.SuccessResponse("设置组件状态成功"))
}