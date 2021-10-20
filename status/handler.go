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

	id := c.Query("id")
	component := c.Query("component")
	if id == "" || component == "" {
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	var status Status
	err := db.Client.Where("id = ? AND component = ?", id, component).Find(&status).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "数据库查询错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(status.Status))
}
func PostStatusHandlerFunc(c *gin.Context) {
	var param struct{
		ID string `json:"id"`
		Component string `json:"component"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	var exist int64
	err := db.Client.Model(&Status{}).Where(
		"id = ? AND component = ?", param.ID, param.Component).Count(&exist).Error
	if err != nil || exist == 0 {
		err := db.Client.Create(&Status{
			ID:      param.ID,
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
			"id = ? AND component = ?", param.ID, param.Component).Update("status", 1).Error
		if err != nil {
			klog.Errorln("数据库错误:", err.Error())
		}
	}

	// 增加修改App的状态
	var total int64
	err = db.Client.Model(&Status{}).Where("status =1 AND id = ?", param.ID).Count(&total).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
	}
	var a app.App
	err = db.Client.First(&a, "id = ?", param.ID).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
	}
	var manifest v1alpha1.Manifest
	err = yaml.Unmarshal([]byte(a.Manifest), &manifest)
	if err != nil {
		klog.Errorln("序列化描述文件错误:", err.Error())
	}
	if int64(len(manifest.Spec.Components)) == total && total != 0 {
		err = db.Client.Model(&app.App{}).Where("status =1 AND id = ?", param.ID).Update("status", 2).Error
		if err != nil {
			klog.Errorln("数据库更新错误:", err.Error())
		}
	}
	c.JSON(200, utils.SuccessResponse("更新成功"))
}