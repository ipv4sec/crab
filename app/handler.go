package app

import (
	"context"
	"crab/aam/v1alpha1"
	"crab/cluster"
	"crab/db"
	"crab/deployment"
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

type DTO struct {
	*App
	Status        string                 `json:"status"`
	Entry         string                 `json:"entry"`
	Dependencies  map[string]interface{} `json:"dependencies"`
	Configurations map[string]interface{} `json:"userconfigs"`
}

type Status struct {
	Name string `json:"name"`
	Message string `json:"message"`
}

func GetAppsHandlerFunc(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	var apps []App
	var total int64
	err := db.Client.Limit(limit).Offset(offset).Where("status > ? AND status < ?", 0, 4).Find(&apps).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库查询错误"))
		return
	}
	err = db.Client.Model(&App{}).Count(&total).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库查询错误"))
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
	var val []DTO
	for i := 0; i < len(apps); i++ {
		var manifest v1alpha1.Application
		err = yaml.Unmarshal([]byte(apps[i].Manifest), &manifest)
		if err != nil {
			klog.Errorln("解析描述文件错误:", err.Error())
			continue
		}

		dependencies := map[string]interface{}{}
		for i := 0; i < len(manifest.Spec.Dependencies); i++ {
			d := Dependency{
				Instances: []Instance{},
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
						d.Instances = append(d.Instances, Instance{ID: apps[j].ID, Version: apps[j].Version})
					}
				}
			}
			dependencies[manifest.Spec.Dependencies[i].Name] = d
		}
		dto := DTO{
			App:    &apps[i],
			Status: "未部署",
			Entry:         endpoints[apps[i].ID+"-http"],
			Configurations: manifest.Spec.Userconfigs,
			Dependencies:  dependencies,
		}
		if apps[i].Status == 1 {
			dto.Status = "正在部署中"
		}
		if apps[i].Status == 2 {
			dto.Status = "部署完成"
		}
		if apps[i].Status == 3 {
			dto.Status = "删除中"
		}
		val = append(val, dto)
	}
	c.JSON(200, utils.SuccessResponse(Pagination{
		Total: len(apps),
		Rows:  val,
	}))
}

func GetAppHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var app App
	err := db.Client.Where("id = ?", id).Find(&app).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseBadRequest, "该实例不存在"))
		return
	}

	island, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(), "island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterInternalServer, "获取根域失败"))
		return
	}
	v, _ := island.Data["root-domain"]
	deploy, err := yaml.Marshal(deployment.Deployment{
		ID:             id,
		Domain:         v,
		Configurations: app.Parameters,
		Dependencies:   app.Dependencies,
	})
	if err != nil {
		klog.Errorln("序列化失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "序列化失败"))
		return
	}
	c.JSON(200, utils.SuccessResponse(map[string]string{
		"id": id,
		"deployment": string(deploy),
	}))
}

func GetAppStatusHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var val []Status
	err := db.Client.Find(&val).Where("id = ?", id).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseBadRequest, "数据库查询错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(val))
}


func PostAppHandlerFunc(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "接收文件错误"))
		return
	}
	currentTimestamp := time.Now().Unix()
	err = c.SaveUploadedFile(file, fmt.Sprintf("/tmp/%v.zip", currentTimestamp))
	if err != nil {
		klog.Errorln("保存文件错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "保存文件错误"))
		return
	}

	err = utils.UnZip(fmt.Sprintf("/tmp/%v", currentTimestamp), fmt.Sprintf("/tmp/%v.zip", currentTimestamp))
	if err != nil {
		klog.Errorln("解压文件错误:", err.Error())
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("/tmp/%v/manifest.yaml", currentTimestamp))
	if err != nil {
		klog.Errorln("读取描述文件错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "读取描述文件错误"))
		return
	}

	var manifest v1alpha1.Application
	err = yaml.Unmarshal(bytes, &manifest)
	if err != nil {
		klog.Errorln("解析描述文件错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "解析描述文件错误"))
		return
	}
	klog.Info("此实例的配置:", manifest.Spec.Userconfigs)
	klog.Info("此实例的依赖:", manifest.Spec.Dependencies)

	configurations := map[string]interface{}{}
	if manifest.Spec.Userconfigs != nil {
		configurations = manifest.Spec.Userconfigs
	}

	configurationsBytes, err := json.Marshal(configurations)
	if err != nil {
		klog.Errorln("序列化运行时配置字段错误:", err.Error())
	}

	dependenciesBytes, err := json.Marshal(manifest.Spec.Dependencies)
	if err != nil {
		klog.Errorln("序列化依赖字段错误:", err.Error())
	}

	app := App{
		ID:     fmt.Sprintf("ins%v", time.Now().Unix()),
		Status: 0,

		Name:          manifest.Metadata.Name,
		Version:       manifest.Metadata.Version,
		Configurations: string(configurationsBytes),
		Dependencies:  string(dependenciesBytes),

		Manifest: string(bytes),

		Parameters: "",
		Deployment: "",
	}
	err = db.Client.Create(&app).Error
	if err != nil {
		klog.Errorln("数据库保存错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库保存错误"))
		return
	}

	dependencies := map[string]interface{}{}
	for i := 0; i < len(manifest.Spec.Dependencies); i++ {
		d := Dependency{
			Instances: []Instance{},
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
					d.Instances = append(d.Instances, Instance{ID: apps[j].ID, Version: apps[j].Version})
				}
			}
		}
		dependencies[manifest.Spec.Dependencies[i].Name] = d
	}

	c.JSON(200, utils.SuccessResponse(struct {
		Dependencies   map[string]interface{} `json:"dependencies" `
		ID             string                 `json:"id"`
		Configurations map[string]interface{} `json:"userconfigs"`
	}{
		Dependencies:   dependencies,
		ID:             app.ID,
		Configurations: configurations,
	}))
}
func PutAppHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var param struct {
		Status         int                     `json:"status"`
		Configurations interface{}             `json:"userconfigs"`
		Dependencies   []struct {
			ID string `json:"id"`
			Name string `json:"name"`
			EntryService string
		} `json:"dependencies"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		klog.Errorln("参数错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var app App
	err = db.Client.Where("id = ?", id).Find(&app).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "该实例不存在"))
		return
	}
	if param.Status == 1 {
		island, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
			Get(context.Background(), "island-info", metav1.GetOptions{})
		if err != nil {
			klog.Errorln("获取根域失败", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrClusterInternalServer, "获取根域失败"))
			return
		}
		v, _ := island.Data["root-domain"]

		for i := 0; i < len(param.Dependencies); i++ {
			if param.Dependencies[i].ID != "" {
				var a App
				err = db.Client.Where("id = ?", param.Dependencies[i].ID).Find(&a).Error
				if err != nil {
					klog.Errorln("数据库查询错误:", err.Error())
					continue
				}
				var manifest v1alpha1.Application
				err = yaml.Unmarshal([]byte(a.Manifest), &manifest)
				if err != nil {
					klog.Errorln("解析描述文件错误:", err.Error())
					continue
				}
				for j := 0; j < len(manifest.Spec.Workloads); j++ {
					if utils.ContainsTrait(manifest.Spec.Workloads[j].Traits, "ingress") {
						param.Dependencies[i].EntryService = manifest.Spec.Workloads[j].Name
					}
				}
			}
		}

		mirror, _ := island.Data["mirror"]
		savedMirrorPath := "/usr/local/workloads/"
		err = utils.InitRepo(savedMirrorPath, mirror)
		if err != nil {
			klog.Errorln("更新工作负载错误:", err.Error())
		}
		val, err := provider.Yaml(app.Manifest, app.ID, v, param.Configurations,
			provider.ConvertToDependency(param.Dependencies), savedMirrorPath)
		if err != nil {
			klog.Errorln("连接到翻译器错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "连接到翻译器错误"))
			return
		}
		err = db.Client.Model(App{}).Where("pk = ?", app.PK).Updates(map[string]interface{}{
			"status": 1, "deployment": val}).Error
		if err != nil {
			klog.Errorln("数据库更新错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "更新状态错误"))
			return
		}

		err = provider.Exec(app.ID, val)
		if err != nil {
			klog.Errorln("调度器执行失败:", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "调度器执行失败"))
			return
		}

		c.JSON(200, utils.SuccessResponse("部署成功"))
		return
	}
	c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
}

func DeleteAppHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	var app App
	err := db.Client.Where("id = ?", id).Find(&app).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "该实例不存在"))
		return
	}

	err = db.Client.Model(App{}).Where("pk = ?", app.PK).Update("status", 3).Error
	if err != nil {
		klog.Errorln("数据库更新错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "更新状态错误"))
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
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "更新状态错误"))
		return
	}
	klog.Infoln("执行命令结果:", output)
	c.JSON(200, utils.SuccessResponse("删除完成"))
}
