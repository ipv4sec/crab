package manifest

import (
	dependency "crab/dependencies"
	"crab/utils"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
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

type contextObj struct {
	AppName       string `json:"appName"`
	ComponentName string `json:"componentName"`
}

type manifestParam struct {
	AppName   string `json:"appName"`
	Namespace string `json:"namespace"`
}

type result struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

type Params struct {
	Content string 	`json:"Content"`
	Instanceid string `json:"InstanceId"`
	Userconfig map[string]string `json:"UserConfig"`
	Dependencies []dependency.Dependency `json:"Dependencies"`
	RootDomain string `json:"RootDomain"`
}

func PostManifestHandlerFunc(c *gin.Context) {
	var err error
	p := Params{}
	err = c.BindJSON(&p)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数格式错误"))
		return
	}
	if p.Content == "" || p.Instanceid == "" || p.RootDomain == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	userconfigStr, err := json.Marshal(p.Userconfig)
	if err != nil {
		c.JSON(200, result{10101,"序列化失败"})
		return
	}

	//生成vale.yaml文件
	vale, err := GenValeYaml(p.Instanceid, p.Content, string(userconfigStr), p.RootDomain, p.Dependencies)
	if err != nil {
		returnData := result{
			10101,
			err.Error(),
		}
		c.JSON(200, returnData)
		return
	}
	//str,err := json.Marshal(vale)
	//if err != nil {
	//	klog.Errorln(err)
	//	return
	//}
	//ioutil.WriteFile("tmp/vela.json", str, 0644)

	//生成k8s.yaml文件
	k8s, err := GenK8sYaml(p.Instanceid, vale)
	if err != nil {
		returnData := result{
			10101,
			err.Error(),
		}
		klog.Errorln(err.Error())
		c.JSON(200, returnData)
		return
	}
	//ioutil.WriteFile("tmp/k8s.yaml", []byte(k8s), 0644)
	returnData := struct {
		Code   int    `json:"code"`
		Result string `json:"result"`
	}{
		0,
		k8s,
	}
	c.JSON(200, returnData)
}

//由manifest.yaml生成vale.yaml
func GenValeYaml(instanceId, content, userconfig, rootDomain string, dependencies []dependency.Dependency) (VelaYaml, error) {
	var vela = VelaYaml{"", make(map[string]interface{}, 0)}
	var err error
	manifestServiceOrigin := ManifestServiceOrigin{}
	err = yaml.Unmarshal([]byte(content), &manifestServiceOrigin)
	if err != nil {
		return vela, err
	}
	vela.Name = manifestServiceOrigin.Metadata.Name

	//components
	if len(manifestServiceOrigin.Spec.Components) == 0 {
		return vela, errors.New("组件不能为空")
	}

	//traits:ingress的组件
	serviceEntryName := entryService(manifestServiceOrigin.Spec.Components)
	if serviceEntryName == "" {
		return vela, errors.New("应用不可访问, 缺少traits.ingress")
	}

	authorizationData, serviceEntryData, configmapData, err := parseDependencies(dependencies)
	if err != nil {
		return vela, err
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

	//configmap
	configItemData := make([]ConfigItemDataItem, 0)
	for k, v := range configmapData {
		configItemData = append(configItemData, ConfigItemDataItem{Name: fmt.Sprintf("%s.host", k), Value: v})
	}
	//添加应用时填写的运行时配置
	configItemData = append(configItemData, ConfigItemDataItem{Name: "userconfig", Value: userconfig})

	for _, svc := range manifestServiceOrigin.Spec.Components {
		service := serviceVela(svc, instanceId, authorizationData, serviceEntryData, configItemData, rootDomain, serviceEntryName)
		vela.Services[svc.Name] = service
	}

	return vela, nil
}

//由vale.yaml生成k8s
func GenK8sYaml(instanceid string, vela VelaYaml) (string, error) {
	//manifest
	manifestK8s, err := GenManifestK8s(instanceid, vela)
	if err != nil {
		klog.Errorln(err.Error())
		return "", err
	}
	//components
	ns := `
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
`
	ns = fmt.Sprintf(ns, instanceid, vela.Name, instanceid)
	componentK8s, err := GenComponentsK8s(vela)
	if err != nil {
		klog.Errorln(err)
		return "", err
	}
	return ns + manifestK8s + componentK8s, nil
}

//获取cue模板
func template(workloadType string) (string, error) {
	var err error
	templatePath := fmt.Sprintf("assets/workloads/%s.cue", workloadType)
	path, _ := filepath.Abs(templatePath)
	if !FileExist(path) {
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
		service := WebserviceVela{
			Workload:      svc.Type,
			Type:          svc.Type,
			Image:         svc.Properties.Image,
			Configs: 	   append(svc.Properties.Configs, ConfigItem{"/etc/configs", "", configItemData}),
			Storage:       svc.Properties.Storage,
			Init:          svc.Properties.Init,
			After:         svc.Properties.After,
			Port:          0,
			Cmd:           svc.Properties.Cmd,
			Args:          svc.Properties.Args,
			Env:           make([]EnvItem, 0),
			Traits:        svc.Traits,
			Authorization: authorization,
			Serviceentry:  serviceentry,
			Namespace:     instanceid,
			Entry:         Entry{},
		}
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
		service := WorkerVela{
			Workload:      svc.Type,
			Type:          svc.Type,
			Image:         svc.Properties.Image,
			Cmd:           svc.Properties.Cmd,
			Args:          svc.Properties.Args,
			Env:           make([]EnvItem, 0),
			After:         svc.Properties.After,
			Init:          svc.Properties.Init,
			Configs:       append(svc.Properties.Configs, ConfigItem{"/etc/configs", "", configItemData}),
			Storage:       svc.Properties.Storage,
			Authorization: authorization,
			Serviceentry:  serviceentry,
			Namespace:     instanceid,
		}
		//service.Configs = append(service.Configs?, ConfigItem{"/etc/configs", "", configItemData})
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
		service := MysqlVela{
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
		service := RedisVela{
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
func parseDependencies(dependencies []dependency.Dependency) ([]dependency.Authorization, []dependency.ServiceEntry, map[string]string, error) {
	var err error
	authorization := make([]dependency.Authorization, 0)
	serviceEntry := make([]dependency.ServiceEntry, 0)
	//dependencies := make([]dependency.Dependency, 0)
	configmap := make(map[string]string, 0)
	//err = json.Unmarshal([]byte(str), &dependencies)
	//if err != nil {
	//	klog.Errorln("依赖解析失败")
	//	return authorization, serviceEntry, configmap, errors.New("依赖解析失败")
	//}
	//解析uses
	dependencyVelas := make([]dependency.DependencyVela, 0)
	for _, v := range dependencies {
		if v.Instanceid != "" && v.EntryService == "" {
			return authorization, serviceEntry, configmap, errors.New("dependencies.entryService不能为空")
		}
		resource,err := dependency.ApiParse(v.Uses)
		if err != nil {
			return authorization, serviceEntry, configmap, err
		}
		dependencyVelas = append(dependencyVelas, dependency.DependencyVela{
			v.Instanceid,
			v.Name,
			v.Location,
			v.Version,
			v.EntryService,
			resource,
		})
	}

	authorization, serviceEntry, configmap, err = dependendService(dependencyVelas)
	if err != nil {
		return authorization, serviceEntry, configmap, err
	}
	return authorization, serviceEntry, configmap, err
}

//依赖的服务,授权
func dependendService(dependencyVelas []dependency.DependencyVela) ([]dependency.Authorization, []dependency.ServiceEntry, map[string]string, error) {
	auth := make([]dependency.Authorization, 0)
	//外部服务调用
	svcEntry := make([]dependency.ServiceEntry, 0)
	//运行时配置
	cm := make(map[string]string, 0)

	for _, v := range dependencyVelas {
		if v.Instanceid != "" {
			auth = append(auth, dependency.Authorization{
				v.Instanceid, v.EntryService, v.Resource,
			})
			cm[v.Name] = fmt.Sprintf("%s.%s.svc.cluster.local.", v.EntryService, v.Instanceid)
		} else {
			if v.Location == "" {
				return auth, svcEntry, cm, errors.New("")
			}
			serviceType, err := inExCheck(v.Location)
			if err != nil {
				return auth, svcEntry, cm, err
			}
			if serviceType == "internal" {
				u, err := url.Parse(v.Location)
				if err != nil {
					return auth, svcEntry, cm, err
				}
				arr := strings.Split(u.Host, ".")
				auth = append(auth, dependency.Authorization{arr[0], arr[1], v.Resource})
			} else {
				arr, err := url.ParseRequestURI(v.Location)
				if err != nil {
					klog.Errorln(err.Error())
					return auth, svcEntry, cm, err
				}
				var protocol string
				if arr.Scheme == "https" {
					protocol = "TLS"
				} else if arr.Scheme == "http" {
					protocol = "http"
				} else {
					klog.Errorln("protocol of the location is not http or https.")
					return auth, svcEntry, cm, errors.New("protocol of the location is not http or https.")
				}
				hostArr := strings.Split(arr.Host, ":")
				var port int
				if len(hostArr) == 1 {
					port = 80
				} else {
					port, err = strconv.Atoi(hostArr[1])
					if err != nil {
						klog.Errorln("转int失败")
						return auth, svcEntry, cm, errors.New("转int失败")
					}
				}
				svcEntry = append(svcEntry,
					dependency.ServiceEntry{arr.Host, port, protocol},
				)
			}
		}
	}
	return auth, svcEntry, cm, nil
}

//traits中包含ingress的组件名称
func entryService(components []Component) string {
	for _, svc := range components {
		for _, v := range svc.Traits {
			if v.Ttype == "ingress" {
				return svc.Name
			}
		}
	}
	return ""
}

//是不是内部服务
func inExCheck(location string) (string, error) {
	u, err := url.Parse(location)
	if err != nil {
		return "", err
	}
	arr := strings.Split(u.Host, ".")
	if arr[len(arr)-1] == "local" {
		return "internal", nil
	}
	return "external", nil
}

func GenManifestK8s(instanceid string, vela VelaYaml) (string, error) {
	manifest := make(map[string]manifestParam, 0)
	manifest["manifest"] = manifestParam{
		vela.Name,
		instanceid,
	}
	manifestStr, err := json.Marshal(manifest)
	if err != nil {
		klog.Errorln("manifestStr json.Marshal 失败")
		return "", errors.New("manifestStr json.Marshal 失败")
	}
	//获取cue模板
	manifestCue, err := template("manifest")
	manifestContent := "\nparameter:%s\n%s"
	manifestContent = fmt.Sprintf(manifestContent, manifestStr, manifestCue)
	fileName := RandomString(manifestContent)
	path := fmt.Sprintf("/tmp/%s.cue", fileName)
	err = ioutil.WriteFile(path, []byte(manifestContent), 0644)
	if err != nil {
		klog.Errorln(err.Error())
		return "", err
	}
	command := fmt.Sprintf("/usr/local/bin/cue export -f %s", path)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorln("执行命令错误", err.Error())
	}
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
	return k8sYaml, nil
}

func GenComponentsK8s(vela VelaYaml) (string, error) {
	k8sYaml := ""
	for k, v := range vela.Services {
		ctxObj := make(map[string]contextObj, 0)
		ctxObj["context"] = contextObj{
			vela.Name,
			k,
		}
		finnnalCueFileContent := "%s\nparameter:%s\n%s"
		ctxObjData, err := json.Marshal(ctxObj)
		if err != nil {
			klog.Errorln("ctxObj 序列化失败")
			return "", errors.New("ctxObj 序列化失败")
		}
		serviceItem, err := json.Marshal(v)
		if err != nil {
			klog.Errorln("vela.Services 序列化失败")
			return "", errors.New("vela.Services 序列化失败")
		}
		workload := ""
		if svc, ok := v.(WebserviceVela); ok {
			workload = svc.Workload
		} else if svc, ok := v.(WorkerVela); ok {
			workload = svc.Workload
		} else if svc, ok := v.(MysqlVela); ok {
			workload = svc.Workload
		} else if svc, ok := v.(RedisVela); ok {
			workload = svc.Workload
		} else {
			klog.Errorln("未知类型的workload")
			return "", errors.New("未知类型的workload")
		}
		template, err := template(workload)
		if err != nil {
			klog.Errorln(err.Error())
			return "", err
		}
		content := fmt.Sprintf(finnnalCueFileContent, ctxObjData, serviceItem, template)
		fileName := RandomString(content)
		path := fmt.Sprintf("/tmp/%s.cue", fileName)
		err = ioutil.WriteFile(path, []byte(content), 0644)
		if err != nil {
			klog.Errorln(err.Error())
			return "", err
		}
		command := fmt.Sprintf("/usr/local/bin/cue export -f %s", path)
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			klog.Errorln("执行命令错误", err.Error())
			return "", err
		}
		var cmdResult struct {
			Parameter map[string]interface{}            `json:"parameter"`
			Outputs   map[string]map[string]interface{} `json:"outputs"`
		}
		err = json.Unmarshal(output, &cmdResult)
		if err != nil {
			klog.Errorln(err.Error())
			return "", err
		}
		for _, out := range cmdResult.Outputs {
			str, err := yaml.Marshal(out)
			if err != nil {
				klog.Errorln(err.Error())
				return "", err
			}
			k8sYaml += fmt.Sprintf("\n---\n#%s\n%s", k, str)
		}
	}
	return k8sYaml, nil
}
