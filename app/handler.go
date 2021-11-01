package app

import (
	"context"
	"crab/aam/v1alpha1"
	"crab/cluster"
	"crab/db"
	dependency "crab/dependencies"
	"crab/exec"
	"crab/provider"
	"crab/utils"
	"encoding/json"
	"fmt"
	"github.com/blang/semver/v4"
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
	Entry  string `json:"entry"`
	Dependencies  map[string]interface{} `json:"dependencies"`
	Configuration map[string]interface{} `json:"userconfig"`
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
	endpoints := map[string]string{"": "未设置域名"}
	for i := 0; i < len(gws.Items); i++ {
		servers := gws.Items[i].Object["spec"].(map[string]interface{})["servers"].([]interface{})
		if len(servers) == 0 {
			continue
		}
		hosts := servers[0].(map[string]interface{})["hosts"].([]interface{})
		if len(hosts) == 0 {
			continue
		}
		endpoints[gws.Items[i].Object["metadata"].(map[string]interface{})["name"].(string)] =
			"http://" + hosts[0].(string)
	}
	vals := []Instance{}
	//|  status   | 意义  |
	//|  ----  | ----  |
	//| 0  | 未部署 |
	//| 1  | 正在部署中  |
	//| 2  | 部署完成  |
	//| 3  | 卸载中  |
	//| 4  | 卸载完成  |
	for i := 0; i < len(apps); i++ {
		var manifest v1alpha1.Manifest
		err = yaml.Unmarshal([]byte(apps[i].Manifest), &manifest)
		if err != nil {
			klog.Errorln("解析描述文件错误:", err.Error())
			continue
		}

		dependencies := map[string]interface{}{}
		for i := 0; i < len(manifest.Spec.Dependencies); i++ {
			d := Dependency{
				Instances: []struct {
					ID      string `json:"instanceid"`
					Version string `json:"version"`
				}{},
			}
			d.Type, d.Link = Link(manifest.Spec.Dependencies[i].Location)
			if d.Type == Mutable {
				var apps []App
				err = db.Client.Where("name = ?", manifest.Spec.Dependencies[i].Name).Find(&apps).Error
				if err != nil {
					klog.Errorln("数据库查询错误:", err.Error())
					continue
				}
				for j := 0; j < len(apps); j++ {
					v, err := semver.Parse(apps[j].Version)
					if err != nil {
						continue
					}
					ra, err := semver.ParseRange(manifest.Spec.Dependencies[i].Version)
					if ra(v) {
						d.Instances = append(d.Instances, struct {
							ID      string `json:"instanceid"`
							Version string `json:"version"`
						}{ID: apps[j].ID, Version: apps[j].Version})
					}
				}
			}
			dependencies[manifest.Spec.Dependencies[i].Name] = d
		}
		ins := Instance{
			App:    &apps[i],
			Status: "未部署",
			// TODO
			Entry:  endpoints[apps[i].ID +"-http"],
			Configuration: manifest.Spec.Configurations,
			Dependencies: dependencies,
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
		klog.Errorln("保存文件错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "保存文件错误"))
		return
	}
	// 增加校验域名
	island, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(), "island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域失败", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "获取根域失败"))
		return
	}
	v, _ := island.Data["root-domain"]
	if v == "" {
		klog.Errorln("获取根域为默认值:", v)
		c.JSON(200, utils.ErrorResponse(10086,  map[string]interface{}{
			"error":"未设置根域, 请设置根域",
			"todo": 1,
		}))
		return
	}

	err = utils.UnZip(fmt.Sprintf("/tmp/%v", currentTimestamp), fmt.Sprintf("/tmp/%v.zip", currentTimestamp))
	if err != nil {
		klog.Errorln("解压文件错误:", err.Error())
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("/tmp/%v/manifest.yaml", currentTimestamp))
	if err != nil {
		klog.Errorln("读取描述文件错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "读取描述文件错误"))
		return
	}

	var manifest v1alpha1.Manifest
	err = yaml.Unmarshal(bytes, &manifest)
	if err != nil {
		klog.Errorln("解析描述文件错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "解析描述文件错误"))
		return
	}
	klog.Info("此实例的配置:", manifest.Spec.Configurations)
	klog.Info("此实例的依赖:", manifest.Spec.Dependencies)

	configuration := map[string]interface{}{}
	if manifest.Spec.Configurations != nil {
		configuration = manifest.Spec.Configurations
	}

	configurationBytes, err := json.Marshal(configuration)
	if err != nil {
		klog.Errorln("序列化运行时配置字段错误:", err.Error())
	}

	dependenciesBytes, err := json.Marshal(manifest.Spec.Dependencies)
	if err != nil {
		klog.Errorln("序列化依赖字段错误:", err.Error())
	}

	app := App{
		ID:     fmt.Sprintf("ins%v", time.Now().Unix()),
		Status:   0,

		Name:     manifest.Metadata.Name,
		Version:  manifest.Metadata.Annotations.Version,
		Configuration: string(configurationBytes),
		Dependencies: string(dependenciesBytes),

		Manifest: string(bytes),

		Parameters: "",
		Deployment: "",
	}
	err = db.Client.Create(&app).Error
	if err != nil {
		klog.Errorln("数据库保存错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "数据库保存错误"))
		return
	}


	dependencies := map[string]interface{}{}
	for i := 0; i < len(manifest.Spec.Dependencies); i++ {
		d := Dependency{
			Instances: []struct {
				ID      string `json:"instanceid"`
				Version string `json:"version"`
			}{},
		}
		d.Type, d.Link = Link(manifest.Spec.Dependencies[i].Location)
		if d.Type == Mutable {
			var apps []App
			err = db.Client.Where("name = ?", manifest.Spec.Dependencies[i].Name).Find(&apps).Error
			if err != nil {
				klog.Errorln("数据库查询错误:", err.Error())
				continue
			}
			for j := 0; j < len(apps); j++ {
				v, err := semver.Parse(apps[j].Version)
				if err != nil {
					continue
				}
				ra, err := semver.ParseRange(manifest.Spec.Dependencies[i].Version)
				if ra(v) {
					d.Instances = append(d.Instances, struct {
						ID      string `json:"instanceid"`
						Version string `json:"version"`
					}{ID: apps[j].ID, Version: apps[j].Version})
				}
			}
		}
		dependencies[manifest.Spec.Dependencies[i].Name] = d
	}

	c.JSON(200, utils.SuccessResponse(struct {
		Dependencies   map[string]interface{} `json:"dependencies" `
		ID             string                 `json:"instanceid"`
		Configurations map[string]interface{}           `json:"userconfig"`
	}{
		Dependencies: dependencies,
		ID:             app.ID,
		Configurations: configuration,
	}))
}
func PutAppHandlerFunc(c *gin.Context) {
	var param struct{
		Status         int         `json:"status"`
		ID             string      `json:"instanceid"` // TODO
		Configurations  interface{} `json:"userconfig"` // TODO
		Dependencies []dependency.Dependency `json:"dependencies"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		klog.Errorln("参数错误", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	if param.ID == "" || param.Status == 0 {
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	// 运行或者卸载
	var app App
	err = db.Client.Where("id = ?", param.ID).Find(&app).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "该实例不存在"))
		return
	}
	// 卸载
	if param.Status == 3 {
		err = db.Client.Model(App{}).Where("pk = ?", app.PK).Update("status", 3).Error
		if err != nil {
			klog.Errorln("数据库更新错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "更新状态错误"))
			return
		}
		executor := exec.CommandExecutor{}
		command := fmt.Sprintf("/usr/local/bin/kubectl delete ns %s", app.ID)
		output, err := executor.ExecuteCommandWithCombinedOutput("bash", "-c", command)
		if err != nil {
			klog.Errorln("执行命令错误", err.Error())
		}
		err = db.Client.Model(App{}).Where(
			"pk = ?", app.PK).Update("status", 4).Error
		if err != nil {
			klog.Errorln("数据库更新错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "更新状态错误"))
			return
		}
		klog.Infoln("执行命令结果:", output)
		c.JSON(200, utils.SuccessResponse("卸载完成"))
		return
	}

	// 运行
	if param.Status == 1 {
		island, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
			Get(context.Background(), "island-info", metav1.GetOptions{})
		if err != nil {
			klog.Errorln("获取根域失败", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "获取根域失败"))
			return
		}
		v, _ := island.Data["root-domain"]

		for i := 0; i < len(param.Dependencies); i++ {
			if param.Dependencies[i].Instanceid != "" {
				var a App
				err = db.Client.Where("id = ?", param.Dependencies[i].Instanceid).Find(&a).Error
				if err != nil {
					klog.Errorln("数据库查询错误:", err.Error())
					continue
				}
				var manifest v1alpha1.Manifest
				err = yaml.Unmarshal([]byte(a.Manifest), &manifest)
				if err != nil {
					klog.Errorln("解析描述文件错误:", err.Error())
					continue
				}
				for j := 0; j < len(manifest.Spec.Components); j++ {
					if utils.ContainsTrait(manifest.Spec.Components[j].Traits, "ingress") {
						param.Dependencies[i].EntryService = manifest.Spec.Components[j].Name
					}
				}
			}
		}

		// TODO
		yaml, err := provider.Yaml(app.Manifest, app.ID, v, param.Configurations, param.Dependencies)
		if err != nil {
			klog.Errorln("连接到翻译器错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "连接到翻译器错误"))
			return
		}
		err = cluster.Client.Apply(context.Background(), []byte(yaml))
		if err != nil {
			klog.Errorln("执行命令错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "执行命令错误"))
			return
		}
		err = db.Client.Model(App{}).Where("pk = ?", app.PK).Updates(map[string]interface{}{
			"status": 1, "deployment": yaml}).Error
		if err != nil {
			klog.Errorln("数据库更新错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(10086, "更新状态错误"))
			return
		}
		c.JSON(200, utils.SuccessResponse("部署成功"))
		return
	}
	c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
}

func DeleteAppHandlerFunc(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, utils.ErrorResponse(10086, "参数错误"))
		return
	}
	err := db.Client.Model(&App{}).Where("id = ?", id).Delete(App{}).Error
	if err != nil {
		klog.Errorln("删除实例错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(10086, "删除实例错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse("删除成功"))
}
