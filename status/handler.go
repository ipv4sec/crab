package status

import (
	"crab/aam/v1alpha1"
	"crab/app"
	"crab/db"
	"crab/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

func GetStatusHandlerFunc(c *gin.Context) {

	name := c.Query("name")
	component := c.Query("component")
	if name == "" || component == "" {
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	var status Status
	err := db.Client.Where("name = ? AND component = ?", name, component).Find(&status).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "数据库查询错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(status.Status))
}
func PostStatusHandlerFunc(c *gin.Context) {
	var param struct{
		Name string `json:"name"`
		Component string `json:"component"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	var exist int64
	err := db.Client.Model(&Status{}).Where(
		"name = ? AND component = ?", param.Name, param.Component).Count(&exist).Error
	if err != nil || exist == 0 {
		err := db.Client.Create(&Status{
			Name:      param.Name,
			Component: param.Component,
			Status:    1,
		}).Error
		if err != nil {
			klog.Errorln("数据库错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
			return
		}
	} else {
		err = db.Client.Model(&Status{}).Where(
			"name = ? AND component = ?", param.Name, param.Component).Update("status", 1).Error
		if err != nil {
			klog.Errorln("数据库错误:", err.Error())
		}
	}

	// 增加修改App的状态
	var total int64
	err = db.Client.Model(&Status{}).Where("status =1 AND name = ?", param.Name).Count(&total).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
	}
	var a app.App
	err = db.Client.First(&a, "uuid = ?", param.Name).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
	}
	var manifest v1alpha1.Manifest
	err = yaml.Unmarshal([]byte(a.Manifest), &manifest)
	if err != nil {
		klog.Errorln("序列化描述文件错误:", err.Error())
	}
	if int64(len(manifest.Spec.Components)) == total && total != 0 {
		err = db.Client.Model(&app.App{}).Where("status =1 AND uuid = ?", param.Name).Update("status", 2).Error
		if err != nil {
			klog.Errorln("数据库更新错误:", err.Error())
		}
	}
	c.JSON(200, utils.SuccessResponse("更新成功"))
}