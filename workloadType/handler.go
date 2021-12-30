package workloadType

import (
	"crab/aam/v1alpha1"
	"crab/db"
	"crab/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
	"strconv"
)

type Pagination struct {
	Total int64         `json:"total"`
	Rows  interface{} `json:"rows"`
}

func GetTypesHandlerFunc(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	var workloadTypes []WorkloadType
	var total int64
	err := db.Client.Limit(limit).Offset(offset).Find(&workloadTypes).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "内部错误"))
		return
	}
	err = db.Client.Model(&WorkloadType{}).Count(&total).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "内部错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(Pagination{
		Total: total,
		Rows:  workloadTypes,
	}))
}

func GetTypeHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var val WorkloadType
	err := db.Client.Where("id = ?", id).Find(&val).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseBadRequest, "该资源不存在"))
		return
	}
	if val.Id == 0 {
		err = db.Client.Where("name = ?", id).Find(&val).Error
		if err != nil {
			klog.Errorln("数据库查询错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseBadRequest, "该资源不存在"))
			return
		}
	}
	if val.Ver == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseBadRequest, "该资源不存在"))
		return
	}
	c.JSON(200, utils.SuccessResponse(val))
}

func PostTypeHandlerFunc(c *gin.Context) {
	var param struct {
		Value string `json:"value"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		klog.Errorln("参数错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var val v1alpha1.WorkloadType
	err = yaml.Unmarshal([]byte(param.Value), &val)
	if err != nil {
		klog.Errorln("反序列化错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	err = db.Client.Model(WorkloadType{}).Create(&WorkloadType{
		Name:      val.Metadata.Name,
		Ver:       val.ApiVersion,
		Value:     param.Value,
		Type:      1,
	}).Error
	if err != nil {
		klog.Errorln("保存到数据库错误", err.Error())
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "名称重复"))
			return
		}
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "内部错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse("创建成功"))
	return
}
func PutTypeHandlerFunc(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var val WorkloadType
	err := db.Client.Where("id = ?", id).Find(&val).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseBadRequest, "该资源不存在"))
		return
	}
	if val.Type == 0 {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "该资源无法删除"))
		return
	}
	var param struct {
		Value string `json:"value"`
	}
	err = c.ShouldBindJSON(&param)
	if err != nil {
		klog.Errorln("参数错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var v1alpha1WorkloadType v1alpha1.WorkloadType
	err = yaml.Unmarshal([]byte(param.Value), &v1alpha1WorkloadType)
	if err != nil {
		klog.Errorln("反序列化错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	err = db.Client.Model(WorkloadType{Id: int64(id)}).Updates(WorkloadType{
		Name:      v1alpha1WorkloadType.Metadata.Name,
		Ver:       v1alpha1WorkloadType.ApiVersion,
		Value:     param.Value,
		Type:      1,
	}).Error
	if err != nil {
		klog.Errorln("保存到数据库错误", err.Error())
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "名称重复"))
			return
		}
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "内部错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse("修改成功"))
	return
}

func DeleteTypeHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var val WorkloadType
	err := db.Client.Where("id = ?", id).Find(&val).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseBadRequest, "该资源不存在"))
		return
	}
	if val.Type == 0 {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "该资源无法删除"))
		return
	}
	err = db.Client.Delete(&WorkloadType{}, id).Error
	if err != nil {
		klog.Errorln("删除错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "删除错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse("删除完成"))
}
