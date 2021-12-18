package deployment

import (
	"crab/aam/v1alpha1"
	"crab/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func PutDeploymentHandlerFunc(c *gin.Context) {
	// id := c.Param("id")
	manifestFileHeader, err := c.FormFile("manifest")
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

	//island, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
	//	Get(context.Background(), "island-info", metav1.GetOptions{})
	//if err != nil {
	//	klog.Errorln("获取根域失败", err.Error())
	//	c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "获取根域失败"))
	//	return
	//}

	//val, err := provider.Yaml(string(manifestBytes), instance.ID, instance.Domain, instance.Configurations,
	//	provider.ConvertToDependency(instance.Dependencies), savedMirrorPath)
	//if err != nil {
	//	klog.Errorln("连接到翻译器错误:", err.Error())
	//	c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "连接到翻译器错误"))
	//	return
	//}
	// TODO
	//err = db.Client.Table("t_app").Where("id = ?", instance.ID).Updates(map[string]interface{}{
	//	"status": 2}).Error
	//if err != nil {
	//	klog.Errorln("数据库更新错误:", err.Error())
	//	c.JSON(200, utils.ErrorResponse(0, "更新状态错误"))
	//	return
	//}

	//err = provider.Exec(instance.ID, val)
	//if err != nil {
	//	klog.Errorln("调度器执行失败:", err.Error())
	//	c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "更新状态错误"))
	//	return
	//}
	c.JSON(200, utils.SuccessResponse("部署成功"))
}