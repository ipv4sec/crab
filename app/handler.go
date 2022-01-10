package app

import (
	"bytes"
	"context"
	"crab/aam/v1alpha1"
	"crab/cluster"
	"crab/db"
	"crab/exec"
	"crab/provider"
	"crab/utils"
	"encoding/json"
	"fmt"
	"github.com/blang/semver/v4"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"strconv"
	"time"
)

var (
	executor = exec.CommandExecutor{}
)

type Pagination struct {
	Total int64         `json:"total"`
	Rows  interface{} `json:"rows"`
}

type DTO struct {
	*App
	Dependencies  map[string]interface{} `json:"dependencies"`
	Configurations map[string]interface{} `json:"userconfigs"`
}

type Logs struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func GetAppsHandlerFunc(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	var apps []App
	var total int64
	err := db.Client.Limit(limit).Offset(offset).Where("status = ?", 1).Find(&apps).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库查询错误"))
		return
	}
	err = db.Client.Model(&App{}).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		klog.Errorln("数据库查询错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库查询错误"))
		return
	}
	for i := 0; i < len(apps); i++ {
		apps[i].Entry = "http://"+apps[i].Entry
	}
	c.JSON(200, utils.SuccessResponse(Pagination{
		Total: total,
		Rows:  apps,
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

	v := map[string] interface{}{}

	cronJob, err := cluster.Client.Clientset.BatchV1beta1().CronJobs(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["cronJob"] = cronJob.Items
	}
	daemonSet, err := cluster.Client.Clientset.AppsV1().DaemonSets(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["daemonSet"] = daemonSet.Items
	}
	deployment, err := cluster.Client.Clientset.AppsV1().Deployments(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["deployment"] = deployment.Items
	}
	job, err := cluster.Client.Clientset.BatchV1().Jobs(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["job"] = job.Items
	}
	pod, err := cluster.Client.Clientset.CoreV1().Pods(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["pod"] = pod.Items
	}
	replicaSet, err := cluster.Client.Clientset.AppsV1().ReplicaSets(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["replicaSet"] = replicaSet.Items
	}
	replicationController, err := cluster.Client.Clientset.CoreV1().ReplicationControllers(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["replicationController"] = replicationController.Items
	}
	statefulSet, err := cluster.Client.Clientset.AppsV1().StatefulSets(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["statefulSet"] = statefulSet.Items
	}
	service, err := cluster.Client.Clientset.CoreV1().Services(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["service"] = service.Items
	}
	configMap, err := cluster.Client.Clientset.CoreV1().ConfigMaps(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["configMap"] = configMap.Items
	}
	pvc, err := cluster.Client.Clientset.CoreV1().PersistentVolumeClaims(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["pvc"] = pvc.Items
	}
	secret, err := cluster.Client.Clientset.CoreV1().Secrets(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["secret"] = secret.Items
	}
	roleBinding, err := cluster.Client.Clientset.RbacV1().RoleBindings(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["roleBinding"] = roleBinding.Items
	}
	role, err := cluster.Client.Clientset.RbacV1().Roles(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["role"] = role.Items
	}
	serviceAccount, err := cluster.Client.Clientset.CoreV1().ServiceAccounts(id).List(context.Background(), metav1.ListOptions{})
	if err == nil {
		v["serviceAccount"] = serviceAccount.Items
	}

	c.JSON(200, utils.SuccessResponse(map[string]interface{}{
		"id": id,
		"deployment": app.Deployment,
		"details": v,
	}))
}


func PostAppHandlerFunc(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		klog.Errorln("接收文件错误", err.Error())
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
	island, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Get(context.Background(), "island-info", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取根域失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterInternalServer, "获取根域失败"))
		return
	}
	v, _ := island.Data["root-domain"]

	id := fmt.Sprintf("ins%v", time.Now().Unix())
	app := App{
		ID:    id ,

		Name:          manifest.Metadata.Name,
		Version:       manifest.Metadata.Version,
		Configurations: string(configurationsBytes),
		Dependencies:  string(dependenciesBytes),

		Manifest: string(bytes),
		Entry: fmt.Sprintf("%s.%s", id, v),

		Additional: "",
		Parameters: "",
		Deployment: "",
		Status: 0,
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
					klog.Errorln("解析实例版本错误:", err.Error())
					c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "解析实例版本错误"))
					return
				}
				ra, err := semver.ParseRange(manifest.Spec.Dependencies[i].Version)
				if err != nil {
					klog.Errorln("解析依赖版本错误:", err.Error())
					c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "解析依赖版本错误"))
					return
				}
				if ra(v) {
					d.Instances = append(d.Instances, Instance{ID: apps[j].ID, Name: apps[j].Name})
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
			Name string `json:"name"`

			ID string `json:"id"`
			Location string `json:"location"`

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

		val, err1, err2 := provider.Yaml(app.Manifest, app.ID, app.Entry, param.Configurations,
			provider.ConvertToDependency(param.Dependencies))
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
		parameters, err := json.Marshal(param.Configurations)
		if err != nil {
			parameters = []byte("")
			klog.Errorln("序列化运行时配置错误", err.Error())
		}
		additional, err := json.Marshal(param.Dependencies)
		if err != nil {
			additional = []byte("")
			klog.Errorln("序列化依赖配置错误", err.Error())
		}

		err = db.Client.Model(App{}).Where("pk = ?", app.PK).Updates(map[string]interface{}{
			"deployment": val, "parameters": string(parameters), "additional": additional, "status": 1}).Error
		if err != nil {
			klog.Errorln("数据库更新错误:", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "更新状态错误"))
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
	executor := exec.CommandExecutor{}
	command := fmt.Sprintf("/usr/local/bin/kubectl delete ns %s", app.ID)
	output, _ := executor.ExecuteCommandWithCombinedOutput("bash", "-c", command)
	klog.Infoln("执行命令结果:", output)
	err = db.Client.Delete(&App{}, app.PK).Error
	if err != nil {
		klog.Errorln("数据库删除错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrDatabaseInternalServer, "数据库删除错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse("删除完成"))
}

func GetPodLogsHandlerFunc(c *gin.Context) {
	id := c.Param("id")
	podName := c.Param("pod")
	if id == "" || podName == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, ""))
		return
	}
	logs, err := GetPodLogs(id, podName)
	if err != nil {
		klog.Errorln("获取POD日志错误:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "获取日志错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(logs))
}

func GetPodLogs(ns, name string) ([]Logs, error) {
	v, err := cluster.Client.Clientset.CoreV1().Pods(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	val := []Logs{}
	for i := 0; i < len(v.Spec.Containers); i++ {
		containerLogs, err := GetContainerLogs(ns, name, v.Spec.Containers[i].Name)
		if err != nil {
			val = append(val, Logs{
				Name:  v.Spec.Containers[i].Name,
				Value: err.Error(),
			})
			continue
		}
		val = append(val, Logs{
			Name:  v.Spec.Containers[i].Name,
			Value: containerLogs,
		})
	}
	return val, nil
}

func GetContainerLogs(ns, pod, container string) (string, error) {
	req := cluster.Client.Clientset.CoreV1().Pods(ns).GetLogs(pod, &v1.PodLogOptions{
		Container: container,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
