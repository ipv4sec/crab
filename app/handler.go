package app

import (
	"context"
	"crab/cluster"
	"crab/db"
	"crab/provider"
	"crab/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"strconv"
	"time"
)

type Pagination struct {
	Total int         `json:"total"`
	Rows  interface{} `json:"rows"`
}

type Instance struct {
	*App
	Status string `json:"status"`
	Entry string `json:"entry"`
}

func GetAppHandlerFunc(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	var apps []App
	var total int64
	err := db.Client.Limit(limit).Offset(offset).Find(&apps).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "数据库查询错误"))
		return
	}
	err = db.Client.Model(&App{}).Count(&total).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "数据库查询错误"))
		return
	}
	gvr := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1beta1",
		Resource: "gateways",
	}
	gws, err := cluster.Client.DynamicClient.Resource(gvr).Namespace("island-system").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln("获取路由资源错误", err.Error())
	}
	endpoints := map[string]string{"":"未设置域名"}
	for i := 0; i < len(gws.Items); i++ {
		servers :=  gws.Items[i].Object["spec"].(map[string]interface{})["servers"].([]interface{})
		if len(servers) == 0 {
			continue
		}
		hosts := servers[0].(map[string]interface{})["hosts"].([]interface{})
		if len(hosts) == 0 {
			continue
		}
		endpoints[gws.Items[i].Object["metadata"].(map[string]interface{})["name"].(string)] = hosts[0].(string)
	}
	var vals []Instance
	for i := 0; i < len(apps); i++ {
		ins := Instance{
			App:    &apps[i],
			Status: "未部署",
			Entry: endpoints[apps[i].Name+"-http"],
		}
		if apps[i].Status == 1 {
			ins.Status = "正在部署中"
		}
		if apps[i].Status == 2 {
			ins.Status = "部署完成"
		}
		if apps[i].Status == 3 {
			ins.Status = "卸载中"
		}
		if apps[i].Status == 4 {
			ins.Status = "卸载完成"
		}
		vals = append(vals, ins)
	}
	c.JSON(200, utils.SuccessResponse(Pagination{
		Total: len(apps),
		Rows:  vals,
	}))
}
func PostAppHandlerFunc(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "接收文件错误"))
		return
	}
	currentTimestamp := time.Now().Unix()
	err = c.SaveUploadedFile(file, fmt.Sprintf("/tmp/%v.zip", currentTimestamp))
	if err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "保存文件错误"))
		return
	}
	err = utils.UnZip(fmt.Sprintf("/tmp/%v", currentTimestamp), fmt.Sprintf("/tmp/%v.zip", currentTimestamp))
	if err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "解压文件错误"))
		return
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("/tmp/%v/manifest.yaml", currentTimestamp))
	if err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "读取描述文件错误"))
		return
	}

	var manifest Manifest
	err = yaml.Unmarshal(bytes, &manifest)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "解析描述文件错误"))
		return
	}
	app := App{
		Name:         manifest.Name,
		Version:      manifest.Version,
		Status:       0,
		Namespace:    fmt.Sprintf("apps-%v", time.Now().Unix()),
		Manifest:     string(bytes),
	}
	err = db.Client.Save(app).Error
	if err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "数据库保存错误"))
		return
	}
	// TODO
	c.JSON(200, utils.SuccessResponse(app))
}
func PutAppHandlerFunc(c *gin.Context) {
	// 运行或者卸载
	var param struct{
		Status int
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	// TODO 卸载
	if param.Status == 3 {

	}

	// 运行
	if param.Status == 1 {
		yaml, err := provider.Yaml("", "", "", nil)
		if err != nil {
			klog.Errorln("连接到翻译器错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "连接到翻译器错误"))
			return
		}
		// TODO
		err = cluster.Client.Apply(context.Background(), []byte(yaml))
		if err != nil {
			klog.Errorln("连接到翻译器错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "连接到翻译器错误"))
			return
		}
		c.JSON(200, utils.SuccessResponse("正在部署中"))
		return
	}

	c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
}

func DeleteAppHandlerFunc(c *gin.Context) {
	id := c.Query("instanceid")
	if id != "" {
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	err := db.Client.Model(&App{}).Where("instance_id = ?", id).Delete(App{}).Error
	if err != nil {
		klog.Errorln("删除实例错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "删除实例错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse("删除成功"))
}
