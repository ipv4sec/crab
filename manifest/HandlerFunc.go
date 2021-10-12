package manifest

import (
	"crab/app"
	"crab/db"
	dependency "crab/dependencies"
	"crab/system"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"k8s.io/klog/v2"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func PostManifestHandlerFunc(c *gin.Context) {
	var err error
	content := c.PostForm("content")
	instanceId := c.PostForm("instanceid")
	userconfig := c.DefaultPostForm("userconfig", "{}")
	dependencies := c.DefaultPostForm("dependencies", "[]")
	if content == "" || instanceId == "" {
		returnData := struct {
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			10101,
			"缺少参数",
		}
		c.JSON(200, returnData)
		return
	}
	//生成vale.yaml文件
	vale, err := GenValeYaml(instanceId, content, dependencies, userconfig)
	if err != nil {
		klog.Errorln(err.Error())
		return
	}

	//生成k8s.yaml文件
	k8s, err := GenK8sYaml(instanceId, vale)
	if err != nil {
		klog.Errorln(err.Error())
		return
	}
	returnData := struct {
		Code   int    `json:"code"`
		Result string `json:"result"`
	}{
		0,
		k8s,
	}
	c.JSON(200, returnData)
}
func PutManifestHandlerFunc(c *gin.Context) {
}

var workloadTypes = []string{"webservice", "worker", "redis", "mysql"}

//workload_vela
type webserviceVela struct {
	Workload      string                     `json:"workload"`
	Image         string                     `json:"image"`
	Configs       []ConfigItem               `json:"configs"`
	Init          string                     `json:"init"`
	After         string                     `json:"after"`
	Port          int                        `json:"port,omitempty"`
	Cmd           []string                   `json:"cmd"`
	Args          []string                   `json:"args,omitempty"`
	Env           []EnvItem                  `json:"env"`
	Traits        []string                   `json:"traits"`
	Authorization []dependency.Authorization `json:"authorization"`
	Serviceentry  []dependency.ServiceEntry  `json:"serviceentry"`
	Namespace     string                     `json:"namespace"`
	Type          string                     `json:"type"`
	Entry         Entry                      `json:"entry"`
}

//workload worker
type workerVela struct {
	Workload      string                     `json:"workload"`
	Image         string                     `json:"image"`
	Cmd           []string                   `json:"cmd"`
	Args          []string                   `json:"args,omitempty"`
	Env           []EnvItem                  `json:"env"`
	After         string                     `json:"after"`
	Init          string                     `json:"init"`
	Configs       []ConfigItem               `json:"configs"`
	Storage       Storage                    `json:"storage"`
	Authorization []dependency.Authorization `json:"authorization"`
	Serviceentry  []dependency.ServiceEntry  `json:"serviceentry"`
	Namespace     string                     `json:"namespace"`
	Type          string                     `json:"type"`
	Entry         Entry                      `json:"entry"`
}

//workload redis
type redisVela struct {
	Workload      string                     `json:"workload"`
	After         string                     `json:"after"`
	Authorization []dependency.Authorization `json:"authorization"`
	Serviceentry  []dependency.ServiceEntry  `json:"serviceentry"`
	Namespace     string                     `json:"namespace"`
	Type          string                     `json:"type"`
	Entry         Entry                      `json:"entry"`
}

//workload mysql
type mysqlVela struct {
	Workload      string                     `json:"workload"`
	Rootpwd       string                     `json:"rootpwd"`
	Storage       Storage                    `json:"storage"`
	Init          string                     `json:"init,omitempty"`
	After         string                     `json:"after"`
	Authorization []dependency.Authorization `json:"authorization"`
	Serviceentry  []dependency.ServiceEntry  `json:"serviceentry"`
	Namespace     string                     `json:"namespace"`
	Type          string                     `json:"type"`
	Entry         Entry                      `json:"entry"`
}

func (svc webserviceVela) workload() string {
	return svc.Workload
}
func (svc workerVela) workload() string {
	return svc.Workload
}
func (svc mysqlVela) workload() string {
	return svc.Workload
}
func (svc redisVela) workload() string {
	return svc.Workload
}

//由manifest.yaml生成vale.yaml
func GenValeYaml(instanceId, str, dependencies, userconfig string) (VelaYaml, error) {
	var err error
	manifestServiceOrigin := ManifestServiceOrigin{}
	err = yaml.Unmarshal([]byte(str), &manifestServiceOrigin)
	if err != nil {
		klog.Errorln(err.Error())
		return VelaYaml{}, nil
	}
	vela := VelaYaml{"", make(map[string]interface{}, 0)}
	vela.Name = manifestServiceOrigin.Metadata.Name

	//parseAnnotations(manifestServiceOrigin.Metadata.Annotations)
	//aa, err := json.Marshal(annotations)
	//if err != nil {
	//	klog.Errorln(err.Error())
	//	return velaYaml{}, nil
	//}

	//组件
	if len(manifestServiceOrigin.Spec.Components) == 0 {
		klog.Errorln("组件不能为空")
		return VelaYaml{}, errors.New("组件不能为空")
	}

	//有ingress的组件
	serviceEntryName := entryService(manifestServiceOrigin.Spec.Components)
	klog.Infoln("ingress ", serviceEntryName)

	//vela格式
	authorizationData, serviceEntryData, configmapData, err := parseDependencies(instanceId, dependencies)
	if err != nil {
		klog.Errorln(err.Error())
		return VelaYaml{}, err
	}

	//为每个 service 创建一个 authorization，授权当前应用下的其他服务有访问的权限
	for _, component := range manifestServiceOrigin.Spec.Components {
		authorizationData = append(authorizationData,
			dependency.Authorization{
				Namespace: instanceId,
				Service:   component.Name,
				Resources: make([]dependency.DependencyUseItem, 0)},
		)
	}

	//根域
	rootDomain, err := system.GetDomain()
	if err != nil {
		klog.Errorln(err.Error())
		return VelaYaml{}, err
	}
	klog.Infoln("GetDomain", rootDomain)

	//configmap
	configItemData := make([]ConfigItemDataItem, 0)
	for k, v := range configmapData {
		configItemData = append(configItemData, ConfigItemDataItem{Name: fmt.Sprintf("%s.host", k), Value: v})
	}
	//添加应用时填写的运行时配置
	if userconfig != "" {
		configItemData = append(configItemData, ConfigItemDataItem{Name: "userconfig", Value: userconfig})
	}

	for _, svc := range manifestServiceOrigin.Spec.Components {
		service := serviceVela(svc, instanceId, authorizationData, serviceEntryData, configItemData, rootDomain, serviceEntryName)
		vela.Services[svc.Name] = service
	}

	return vela, nil
}

//由vale.yaml生成k8s
func GenK8sYaml(instanceid string, vela VelaYaml) (string, error) {
	var err error
	var finalYaml string

	//数据库中查找数据
	appInfo := app.App{}
	err = db.Client.Where("id = ?", instanceid).First(&appInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			klog.Errorln("err")
			klog.Errorln("Error: Dependence app not exist. instanceid=" + instanceid)
			return "", errors.New("Error: Dependence app not exist. instanceid=" + instanceid)
		}
		klog.Errorln(err.Error())
		return "", errors.New("查询数据库错误")
	}
	manifestName := vela.Name
	manifest := make(map[string]struct {
		AppName   string `json:"appName"`
		Namespace string `json:"namespace"`
	}, 0)
	manifest["manifest"] = struct {
		AppName   string `json:"appName"`
		Namespace string `json:"namespace"`
	}{
		manifestName,
		instanceid,
	}
	manifestStr, err := json.Marshal(manifest)
	if err != nil {
		klog.Errorln("manifestStr json.Marshal 失败")
		return "", errors.New("manifestStr json.Marshal 失败")
	}
	//生成兜底规则
	manifestCue, err := template("manifest")

	manifestContent := `
parameter:%s
%s
`
	manifestContent = fmt.Sprintf(manifestContent, manifestStr, manifestCue)
	fileName := RandomString(manifestContent)
	path := fmt.Sprintf("/tmp/%s.cue", fileName)
	//path = "tmp/cue/manifest.cue"
	err = ioutil.WriteFile(path, []byte(manifestContent), 0644)
	if err != nil {
		klog.Errorln(err.Error())
		return "", err
	}
	command := fmt.Sprintf("/usr/local/bin/cue export -f %s", path)
	//command = fmt.Sprintf("cue export -f %s", path)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorln("执行命令错误", err.Error())
	}
	//klog.Errorln(string(output))
	var out struct {
		Parameter map[string]interface{}            `json:"parameter"`
		Outputs   map[string]map[string]interface{} `json:"outputs"`
	}
	err = json.Unmarshal(output, &out)
	if err != nil {
		klog.Errorln(err.Error())
		return "", err
	}
	k8sYaml := ""
	for _, output := range out.Outputs {
		str, err := yaml.Marshal(output)
		if err != nil {
			klog.Errorln(err.Error())
			return "", err
		}
		k8sYaml += fmt.Sprintf("---\n#manifest\n%s", str)
	}
	//service
	for k, v := range vela.Services {
		type contextObj struct {
			AppName       string `json:"appName"`
			ComponentName string `json:"componentName"`
		}
		ctxObj := make(map[string]contextObj, 0)
		ctxObj["context"] = contextObj{manifestName, k}
		if err != nil {
			klog.Errorln(err.Error())
			return "", err
		}
		finnnalCueFileContent := `
%s
parameter:%s
%s
`
		ctxObjData, err := json.Marshal(ctxObj)
		if err != nil {
			klog.Errorln("ctxObj json.Marshal 失败")
			return "", errors.New("ctxObj json.Marshal 失败")
		}

		serviceItem, err := json.Marshal(v)
		if err != nil {
			klog.Errorln("vela.Services json.Marshal 失败")
			return "", errors.New("vela.Services json.Marshal 失败")
		}
		workload := ""
		if svc, ok := v.(webserviceVela); ok {
			workload = svc.Workload
		} else if svc, ok := v.(workerVela); ok {
			workload = svc.Workload
		} else if svc, ok := v.(mysqlVela); ok {
			workload = svc.Workload
		} else if svc, ok := v.(redisVela); ok {
			workload = svc.Workload
		} else {
			klog.Errorln("未知1类型的workload")
			return "", errors.New("未知类型的workload")
		}
		template, err := template(workload)
		if err != nil {
			klog.Errorln(err.Error())
			return "", err
		}
		content := fmt.Sprintf(finnnalCueFileContent, ctxObjData, serviceItem, template)
		fileName = RandomString(content)
		path := fmt.Sprintf("/tmp/%s.cue", fileName)
		//path = "tmp/cue/" + k + ".cue"
		err = ioutil.WriteFile(path, []byte(content), 0644)
		if err != nil {
			klog.Errorln(err.Error())
			return "", err
		}
		command = fmt.Sprintf("/usr/local/bin/cue export -f %s", path)
		//command = fmt.Sprintf("cue export -f %s", path)
		cmd = exec.Command("bash", "-c", command)
		output, err = cmd.CombinedOutput()
		if err != nil {
			klog.Errorln("执行命令错误", err.Error())
			return "", err
		}
		err = json.Unmarshal(output, &out)
		if err != nil {
			klog.Errorln(err.Error())
			return "", err
		}
		for _, out := range out.Outputs {
			str, err := yaml.Marshal(out)
			if err != nil {
				klog.Errorln(err.Error())
				return "", err
			}
			k8sYaml += fmt.Sprintf("\n---\n#%s\n%s", k, str)
		}
	}
	k8s := `
apiVersion: v1
kind: Namespace
metadata:
 name: %s
 labels:
   istio-injection: enabled
---
apiVersion: v1
kind: ServiceAccount
metadata:
 name: %s
 namespace: %s
%s
`
	finalYaml = fmt.Sprintf(k8s, instanceid, vela.Name, instanceid, k8sYaml)
	//err = ioutil.WriteFile("tmp/k8s.yaml", []byte(finalYaml), 0644)
	//if err != nil {
	//	klog.Errorln("k8s写入文件失败")
	//	return "", errors.New("k8s写入文件失败")
	//}
	//klog.Errorln("k8s写入文件成功")
	return finalYaml, nil
}

//获取cue模板
func template(workloadType string) (string, error) {
	var err error
	templatePath := fmt.Sprintf("assets/workloads/%s.cue", workloadType)
	path, _ := filepath.Abs(templatePath)
	if ! FileExist(path) {
		klog.Errorln(err.Error())
		return "", errors.New(fmt.Sprintf("文件：%s 不存在", path))
	}
	t, err := ioutil.ReadFile(templatePath)
	if err != nil {
		klog.Errorln(err.Error())
		return "", err
	}
	content := string(t)

	//替换import为真实内容
	re, _ := regexp.Compile("import\\s*\"([^\"]*)\"")
	matchResult := re.FindAllStringSubmatch(content, -1)
	for _, v := range matchResult {
		if len(matchResult) > 0 {
			includeMod, err := template(v[1])
			if err != nil {
				klog.Errorln(err.Error())
				return "", err
			}
			content = strings.ReplaceAll(content, v[0], includeMod)
		}
	}
	//ioutil.WriteFile("2.cue", []byte(content),0644)
	return content, nil
}
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//生成随机字符串
func RandomString(str string) string {
	h := md5.New()
	h.Write([]byte(str + string(rand.Intn(1000))))
	return hex.EncodeToString(h.Sum(nil))
}

//生成kubevela格式的service
func serviceVela(svc Component, instanceid string, authorization []dependency.Authorization, serviceentry []dependency.ServiceEntry, configItemData []ConfigItemDataItem, rootDomain string, serviceEntryName string) interface{} {
	if svc.Type == "webservice" {
		service := webserviceVela{
			Workload:      svc.Type,
			Type:          svc.Type,
			Image:         svc.Properties.Image,
			Configs:       make([]ConfigItem, 0),
			Init:          svc.Properties.Init,
			After:         svc.Properties.After,
			Port:          0,
			Cmd:           svc.Properties.Cmd,
			Args:          svc.Properties.Args,
			Env:           make([]EnvItem, 0),
			Traits:        svc.Properties.Traits,
			Authorization: authorization,
			Serviceentry:  serviceentry,
			Namespace:     instanceid,
			Entry:         Entry{},
		}
		service.Configs = append(service.Configs, ConfigItem{"/etc/configs", "", configItemData})
		if serviceEntryName == svc.Name {
			path := make([]string, 0)
			path = append(path, "/*")
			service.Entry = Entry{
				fmt.Sprintf("%s.%s", instanceid, rootDomain),
				path,
			}
		} else {
			service.Entry = Entry{
				"",
				make([]string, 0),
			}
		}
		return service
	} else if svc.Type == "worker" {
		service := workerVela{
			Workload:      svc.Type,
			Type:          svc.Type,
			Image:         svc.Properties.Image,
			Cmd:           svc.Properties.Cmd,
			Args:          svc.Properties.Args,
			Env:           make([]EnvItem, 0),
			After:         svc.Properties.After,
			Init:          svc.Properties.Init,
			Configs:       make([]ConfigItem, 0),
			Storage:       svc.Properties.Storage,
			Authorization: authorization,
			Serviceentry:  serviceentry,
			Namespace:     instanceid,
		}
		service.Configs = append(service.Configs, ConfigItem{"/etc/configs", "", configItemData})
		if serviceEntryName == svc.Name {
			path := make([]string, 0)
			path = append(path, "/*")
			service.Entry = Entry{
				fmt.Sprintf("%s.%s", instanceid, rootDomain),
				path,
			}
		} else {
			service.Entry = Entry{
				"",
				make([]string, 0),
			}
		}
		return service
	} else if svc.Type == "mysql" {
		service := mysqlVela{
			Workload:      svc.Type,
			Type:          svc.Type,
			Rootpwd:       svc.Properties.Rootpwd,
			Storage:       svc.Properties.Storage,
			Init:          svc.Properties.Init,
			After:         svc.Properties.After,
			Authorization: authorization,
			Serviceentry:  serviceentry,
			Namespace:     instanceid,
		}
		return service
	} else if svc.Type == "redis" {
		service := redisVela{
			Workload:      svc.Type,
			Type:          svc.Type,
			After:         svc.Properties.After,
			Authorization: authorization,
			Serviceentry:  serviceentry,
			Namespace:     instanceid,
		}
		return service
	}
	return nil
}

//处理依赖
func parseDependencies(instanceId, str string) ([]dependency.Authorization, []dependency.ServiceEntry, map[string]string, error) {
	var err error
	authorization := make([]dependency.Authorization, 0)
	serviceEntry := make([]dependency.ServiceEntry, 0)
	dependencies := make([]dependency.Dependency, 0)
	configmap := make(map[string]string, 0)
	err = json.Unmarshal([]byte(str), &dependencies)
	if err != nil {
		klog.Errorln("依赖解析错误")
		return authorization, serviceEntry, configmap, errors.New("依赖解析错误")
	}
	//解析uses
	dependencyVelas := make([]dependency.DependencyVela, 0)
	for _, v := range dependencies {
		dependencyVelas = append(dependencyVelas, dependency.DependencyVela{
			v.Instanceid,
			v.Name,
			v.Location,
			v.Version,
			dependency.ApiParse(v.Uses),
		})
	}

	authorization, serviceEntry, configmap, err = dependendService(instanceId, dependencyVelas)
	if err != nil {
		klog.Errorln(err.Error())
		return authorization, serviceEntry, configmap, err
	}
	return authorization, serviceEntry, configmap, err
}

//依赖的服务,授权
func dependendService(instanceId string, dependencyVelas []dependency.DependencyVela) ([]dependency.Authorization, []dependency.ServiceEntry, map[string]string, error) {
	dependenceAuthorization := make([]dependency.Authorization, 0)
	//外部服务调用
	externalService := make([]dependency.ServiceEntry, 0)
	//运行时配置
	configmap := make(map[string]string, 0)

	for _, v := range dependencyVelas {
		if v.Instanceid != "" { //有Instanceid，说明是内部服务
			appInfo, err := appInfo(v.Instanceid)
			if err != nil {
				klog.Errorln(err.Error())
				return dependenceAuthorization, externalService, configmap, err
			}
			manifestServiceOrigin := ManifestServiceOrigin{}
			err = yaml.Unmarshal([]byte(appInfo.Manifest), &manifestServiceOrigin)
			if err != nil {
				klog.Errorln(err.Error())
				return dependenceAuthorization, externalService, configmap, err
			}
			entryService := entryService(manifestServiceOrigin.Spec.Components)
			dependenceAuthorization = append(dependenceAuthorization, dependency.Authorization{
				v.Instanceid, entryService, v.Resource,
			})
			configmap[appInfo.Name] = fmt.Sprintf("%s.%s.svc.cluster.local.", entryService, v.Instanceid)
		} else {
			if v.Location == "" {
				klog.Errorln("Error: location is empty")
				return dependenceAuthorization, externalService, configmap, errors.New("location is empty")
			}
			if inExCheck(v.Location) == "internal" {
				u, err := url.Parse(v.Location)
				if err != nil {
					klog.Errorln(err.Error())
					return dependenceAuthorization, externalService, configmap, err
				}
				arr := strings.Split(u.Host, ".")
				dependenceAuthorization = append(dependenceAuthorization, dependency.Authorization{arr[0], arr[1], v.Resource})
			} else {
				arr, err := url.ParseRequestURI(v.Location)
				if err != nil {
					klog.Errorln(err.Error())
					return dependenceAuthorization, externalService, configmap, err
				}
				var protocol string
				if arr.Scheme == "https" {
					protocol = "TLS"
				} else if arr.Scheme == "http" {
					protocol = "http"
				} else {
					klog.Errorln("Error: protocol of the location is not http or https.")
					return dependenceAuthorization, externalService, configmap, errors.New("protocol of the location is not http or https.")
				}
				arr2 := strings.Split(arr.Host, ":")
				var port int
				if len(arr2) == 1 {
					port = 80
				} else {
					port, err = strconv.Atoi(arr2[1])
					if err != nil {
						klog.Errorln("转int失败")
						return dependenceAuthorization, externalService, configmap, errors.New("转int失败")
					}
				}
				externalService = append(externalService,
					dependency.ServiceEntry{arr.Host, port, protocol},
				)
			}
		}
	}
	return dependenceAuthorization, externalService, configmap, nil
}

//获取app信息
func appInfo(id string) (app.App, error) {
	var err error
	app := app.App{}
	err = db.Client.Where("id = ?", id).First(&app).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			klog.Errorln("err")
			klog.Errorln("Error: Dependence app not exist. instanceid=" + id)
			return app, errors.New("Error: Dependence app not exist. instanceid=" + id)
		}
		klog.Errorln(err.Error())
		return app, errors.New("查询数据库错误")
	}
	return app, nil
}

//返回traits中包含ingress的服务名称
func entryService(components []Component) string {
	for _, svc := range components {
		klog.Errorln(svc.Properties.Traits)
		for _, v := range svc.Properties.Traits {
			if v == "ingress" {
				return svc.Name
			}
		}
	}
	return ""
}

//是不是内部服务
func inExCheck(location string) string {
	u, err := url.Parse(location)
	if err != nil {
		panic(err)
	}
	arr := strings.Split(u.Host, ".")
	if arr[len(arr)-1] == "local" {
		return "internal"
	} else {
		return "external"
	}
}
