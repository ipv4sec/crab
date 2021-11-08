package parser

import (
	"bytes"
	"crab/aam/v1alpha1"
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
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Params struct {
	Content      string       `json:"Content"`
	Instanceid   string       `json:"InstanceId"`
	Userconfig   interface{}  `json:"UserConfig"`
	Dependencies []Dependency `json:"Dependencies"`
	RootDomain   string       `json:"RootDomain"`
	WorkloadPath string       `json:"WorkloadPath"`
}

func PostManifestHandlerFunc(c *gin.Context) {
	var err error
	p := Params{}
	err = c.BindJSON(&p)
	if err != nil {
		c.JSON(200, Result{ErrBadRequest, "参数错误"})
		return
	}
	if p.Content == "" || p.Instanceid == "" || p.RootDomain == "" || p.WorkloadPath == "" {
		c.JSON(200, Result{ErrBadRequest, "缺少参数"})
		return
	}
	userconfigStr, err := json.Marshal(p.Userconfig)
	if err != nil {
		c.JSON(200, Result{ErrInternalServer, "序列化失败"})
		return
	}
	//解析描述文件
	var application v1alpha1.Application
	err = yaml.Unmarshal([]byte(p.Content), &application)
	if err != nil {
		c.JSON(200, Result{ErrBadRequest, "文件解析失败"})
		return
	}

	//验证参数，返回参数json,返回vendor内容
	workloadResource, err := checkParams(application, p.WorkloadPath)
	if err != nil {
		c.JSON(200, Result{ErrBadRequest, err.Error()})
		return
	}

	//生成vale.yaml文件
	vale, err := GenValeYaml(p.Instanceid, application, string(userconfigStr), p.RootDomain, p.Dependencies)
	if err != nil {
		c.JSON(200, Result{ErrInternalServer, err.Error()})
		return
	}
	//str, err := json.Marshal(vale)
	//if err != nil {
	//	klog.Errorln(err)
	//	return
	//}
	//ioutil.WriteFile("tmp/vela.json", str, 0644)

	//生成k8s.yaml文件
	k8s, err := GenK8sYaml(p.Instanceid, vale, workloadResource)
	if err != nil {
		klog.Errorln(err.Error())
		c.JSON(200, Result{ErrInternalServer, err.Error()})
		return
	}
	k8s2, err := yaml.Marshal(k8s)
	if err != nil {
		fmt.Println(err)
		return
	}
	//ioutil.WriteFile("tmp/k8s.yaml", k8s2, 0644)
	c.JSON(200, Result{0, string(k8s2)})
}

//由manifest.yaml生成vale.yaml
func GenValeYaml(instanceId string, application v1alpha1.Application, userconfig string, rootDomain string, dependencies []Dependency) (VelaYaml, error) {
	var vela = VelaYaml{"", make(map[string]interface{}, 0)}
	var err error
	vela.Name = application.Metadata.Name

	//traits:ingress的组件
	serviceEntryName := entryService(application.Spec.Workloads)
	authorizationData, serviceEntryData, configmapData, err := parseDependencies(dependencies)
	if err != nil {
		return vela, err
	}

	//为每个 service 创建一个 authorization，授权当前应用下的其他服务有访问的权限
	for _, workload := range application.Spec.Workloads {
		authorizationData = append(authorizationData,
			Authorization{
				Namespace: instanceId,
				Service:   workload.Name,
				Resources: make([]DependencyUseItem, 0)},
		)
	}

	//configmap
	configItemData := make([]ConfigItemDataItem, 0)
	for k, v := range configmapData {
		configItemData = append(configItemData, ConfigItemDataItem{Name: fmt.Sprintf("%s.host", k), Value: v})
	}
	//添加应用时填写的运行时配置
	configItemData = append(configItemData, ConfigItemDataItem{Name: "userconfig", Value: userconfig})
	for _, workload := range application.Spec.Workloads {
		service := serviceVela(workload, instanceId, authorizationData, serviceEntryData, configItemData, rootDomain, serviceEntryName)
		vela.Services[workload.Name] = service
	}
	return vela, nil
}

//由vale.yaml生成k8s
func GenK8sYaml(instanceid string, vela VelaYaml, workloadParam map[string]WorkloadParam) (ParserData, error) {
	parserData := ParserData{
		Init:      "",
		Name:      "",
		Workloads: make(map[string]Workload, 0),
	}
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
---
apiVersion: "security.istio.io/v1beta1"
kind: "AuthorizationPolicy"
metadata:
 name: %s
 namespace: %s
spec:
 {}
`
	ns = fmt.Sprintf(ns, instanceid, vela.Name, instanceid, instanceid, instanceid)
	parserData.Init = ns
	parserData.Name = vela.Name
	//处理workload
	for k, v := range vela.Services {
		//ctxObj := make(map[string]ContextObj, 0)
		//ctxObj["context"] = ContextObj{
		//	vela.Name,
		//	k,
		//	instanceid,
		//}
		//finnnalCueFileContent := "%s\nparameter:%s\n%s"
		//template := workloadParams[k].VendorCue
		//ctxObjData, err := json.Marshal(ctxObj)
		//if err != nil {
		//	klog.Errorln("ctxObj 序列化失败")
		//	return parserData, errors.New("ctxObj 序列化失败")
		//}
		//serviceItem, err := json.Marshal(v)
		//if err != nil {
		//	klog.Errorln("vela.Services 序列化失败")
		//	return parserData, errors.New("vela.Services 序列化失败")
		//}
		//content := fmt.Sprintf(finnnalCueFileContent, ctxObjData, serviceItem, template)
		//fileName := RandomString(content)
		//path := fmt.Sprintf("/tmp/test%s.cue", fileName)
		//err = ioutil.WriteFile(path, []byte(content), 0644)
		//if err != nil {
		//	klog.Errorln(err.Error())
		//	return parserData, err
		//}
		//command := fmt.Sprintf("/usr/local/bin/cue export -f %s", path)
		//cmd := exec.Command("bash", "-c", command)
		//output, err := cmd.CombinedOutput()
		//if err != nil {
		//	klog.Errorln("执行命令错误", err.Error())
		//	return parserData, err
		//}
		//var cmdResult map[string]map[string]interface{}
		//err = json.Unmarshal(output, &cmdResult)
		//if err != nil {
		//	klog.Errorln(err.Error())
		//	return parserData, err
		//}
		ctxObj := make(map[string]ContextObj, 0)
		ctxObj["context"] = ContextObj{
			vela.Name,
			k,
			instanceid,
		}
		cmdResult,err := GenWorkloadCue(ctxObj, workloadParam[k], v)
		if err != nil {
			klog.Errorln(err)
			return parserData, err
		}
		var workload Workload
		construct := make(map[string]string, 0)
		count := 0
		for k, out := range cmdResult["construct"] {
			str, err := yaml.Marshal(out)
			if err != nil {
				klog.Errorln(err.Error())
				return parserData, err
			}
			construct[k] = string(str)
			count++
		}
		if count == 0 {
			err = errors.New("vendor未实现type")
			return parserData, err
		}
		workload.Construct = construct
		traits := make(map[string]string, 0)
		if len(workloadParam[k].Traits) > 0 { //有trait
			for _, v := range workloadParam[k].Traits {
				count = 0
				arr := strings.Split(v, "/")
				v = arr[len(arr)-1]
				for k, out := range cmdResult[v] {
					str, err := yaml.Marshal(out)
					if err != nil {
						klog.Errorln(err.Error())
						return parserData, err
					}
					traits[k] = string(str)
					count++
				}
				if count == 0 {
					err = errors.New("未实现trait")
					return parserData, err
				}
			}

		}
		workload.Traits = traits
		parameterStr, err := yaml.Marshal(cmdResult["parameter"])
		if err != nil {
			fmt.Println(err.Error())
			return ParserData{}, nil
		}
		workload.Parameter = string(parameterStr)
		parserData.Workloads[k] = workload
	}
	return parserData, nil
}

func GenWorkloadCue(ctxObj map[string]ContextObj, workloadParam WorkloadParam, workload interface{}) (map[string]map[string]interface{}, error){
	var cmdResult map[string]map[string]interface{}
	finnnalCueFileContent := "%s\nparameter:%s\n%s"
	template := workloadParam.VendorCue
	ctxObjData, err := json.Marshal(ctxObj)
	if err != nil {
		klog.Errorln("ctxObj 序列化失败")
		return cmdResult, errors.New("ctxObj 序列化失败")
	}
	serviceItem, err := json.Marshal(workload)
	if err != nil {
		klog.Errorln("vela.Services 序列化失败")
		return cmdResult, errors.New("vela.Services 序列化失败")
	}
	content := fmt.Sprintf(finnnalCueFileContent, ctxObjData, serviceItem, template)
	fileName := RandomString(content)
	path := fmt.Sprintf("/tmp/test%s.cue", fileName)
	err = ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		klog.Errorln(err.Error())
		return cmdResult, err
	}
	command := fmt.Sprintf("/usr/local/bin/cue export -f %s", path)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorln("执行命令错误", err.Error())
		return cmdResult, err
	}
	err = json.Unmarshal(output, &cmdResult)
	if err != nil {
		klog.Errorln(err.Error())
		return cmdResult, err
	}
	return cmdResult, nil
}

//获取cue模板
func modTemplate(workloadVendor, mod, vendorDir string) (string, error) {
	var err error
	pos := strings.LastIndex(workloadVendor, "/")
	path := fmt.Sprintf("%s/%s/workloadVendor/%s.cue", vendorDir, workloadVendor[:pos+1], mod)
	if !FileExist(path) {
		return "", errors.New(fmt.Sprintf("文件：%s 不存在", path))
	}
	t, err := ioutil.ReadFile(path)
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
			includeMod, err := modTemplate(workloadVendor, v[1], vendorDir)
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
	h.Write([]byte(str + string(rand.Intn(999999))))
	return hex.EncodeToString(h.Sum(nil))
}

//生成kubevela格式的service
func serviceVela(workload v1alpha1.Workload, instanceid string, authorization []Authorization, serviceentry []ServiceEntry, configItemData []ConfigItemDataItem, rootDomain string, serviceEntryName string) interface{} {
	properties := GetProperties(workload.Properties)
	properties["authorization"] = authorization
	properties["serviceentry"] = serviceentry
	configs2 := make([]interface{}, 0)
	if configs, ok := properties["configs"]; ok {
		for _, v := range configs.([]interface{}) {
			configs2 = append(configs2, v)
		}
	}
	configs2 = append(configs2, ConfigItem{"/etc/configs", "", configItemData})
	properties["configs"] = configs2
	if serviceEntryName == workload.Name {
		path := make([]string, 0)
		path = append(path, "/*")
		entry := Entry{
			fmt.Sprintf("%s.%s", instanceid, rootDomain),
			path,
		}
		properties["entry"] = entry
	}
	return properties
}

//处理依赖
func parseDependencies(dependencies []Dependency) ([]Authorization, []ServiceEntry, map[string]string, error) {
	var err error
	authorization := make([]Authorization, 0)
	serviceEntry := make([]ServiceEntry, 0)
	configmap := make(map[string]string, 0)
	//解析uses
	dependencyVelas := make([]DependencyVela, 0)
	for _, v := range dependencies {
		if v.Instanceid != "" && v.EntryService == "" {
			return authorization, serviceEntry, configmap, errors.New("dependencies.entryService不能为空")
		}
		resource, err := ApiParse(v.Items)
		if err != nil {
			return authorization, serviceEntry, configmap, err
		}
		dependencyVelas = append(dependencyVelas, DependencyVela{
			v.Instanceid,
			v.Name,
			v.Location,
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
func dependendService(dependencyVelas []DependencyVela) ([]Authorization, []ServiceEntry, map[string]string, error) {
	auth := make([]Authorization, 0)
	//外部服务调用
	svcEntry := make([]ServiceEntry, 0)
	//运行时配置
	cm := make(map[string]string, 0)

	for _, v := range dependencyVelas {
		if v.Instanceid != "" {
			auth = append(auth, Authorization{
				v.Instanceid, v.EntryService, v.Resource,
			})
			cm[v.Name] = fmt.Sprintf("%s.%s.svc.cluster.local", v.EntryService, v.Instanceid)
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
				auth = append(auth, Authorization{arr[0], arr[1], v.Resource})
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
				svcEntry = append(svcEntry, ServiceEntry{arr.Host, port, protocol})
			}
		}
	}
	return auth, svcEntry, cm, nil
}

//traits中包含ingress的组件名称
func entryService(workloads []v1alpha1.Workload) string {
	for _, svc := range workloads {
		for _, v := range svc.Traits {
			arr := strings.Split(v.Type, "/")
			trait := arr[len(arr)-1]
			if trait == "ingress" {
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

func ApiParse(uses map[string][]string) ([]DependencyUseItem, error) {
	var err error
	rtn := make([]DependencyUseItem, 0)
	for k, v := range uses {
		count := 0
		actions := make([]string, 0)
		for _, option := range v {
			if option == "create" {
				actions = append(actions, "POST")
			} else if option == "read" {
				actions = append(actions, "GET", "HEAD", "OPTIONS")
			} else if option == "update" {
				actions = append(actions, "PUT", "PATCH")
			} else if option == "delete" {
				actions = append(actions, "DELETE")
			} else {
				return rtn, errors.New(fmt.Sprintf("依赖资源的操作类型(%s)不存在\n", option))
			}
			count++
		}
		if count == 0 {
			return rtn, errors.New("依赖资源的操作类型不能为空")
		}
		rtn = append(rtn, DependencyUseItem{k, actions})
	}
	return rtn, err
}

func checkParams(application v1alpha1.Application, vendorDir string) (map[string]WorkloadParam, error) {
	var err error
	returnData := make(map[string]WorkloadParam, 0)
	if len(application.Spec.Workloads) == 0 {
		err = errors.New("application.Spec.Workloads 不能为空")
		return returnData, err
	}
	for _, workload := range application.Spec.Workloads {
		//检查type参数
		err = CheckTypeParam(workload, vendorDir)
		if err != nil {
			klog.Errorln(err)
			return returnData, err
		}
		//检查trait参数
		if len(workload.Traits) > 0 {
			for _, trait := range workload.Traits {
				err = CheckTraitParam(trait, vendorDir)
				if err != nil {
					return returnData, err
				}
			}
		}
		var workloadParams WorkloadParam
		workloadParams.Type = workload.Type
		workloadParams.Vendor = workload.Vendor

		properties := GetProperties(workload.Properties)
		workloadParams.Parameter = properties

		t, _ := GetWorkloadType(workload.Type, vendorDir)
		workloadParams.Traits = t.Spec.Traits

		var v v1alpha1.WorkloadVendor
		v, err = GetWorkloadVendor(workload.Vendor, vendorDir)
		if err != nil {
			return returnData, err
		}
		workloadParams.VendorCue = v.Spec
		returnData[workload.Name] = workloadParams
	}
	return returnData, nil
}

//获取WorkloadType
func GetWorkloadType(typeName, vendorDir string) (v1alpha1.WorkloadType, error) {
	var err error
	var t v1alpha1.WorkloadType
	pos := strings.LastIndex(typeName, "/")
	path := fmt.Sprintf("%s/%s/workloadType/%s.yaml", vendorDir, typeName[:pos+1], typeName[pos+1:])
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("workload.Type: %s 不存在\n", typeName))
		return t, err
	}
	//解析为结构体
	err = yaml.Unmarshal(content, &t)

	return t, err
}

//获取WorkloadVendor
func GetWorkloadVendor(vendorName, vendorDir string) (v1alpha1.WorkloadVendor, error) {
	var err error
	var v v1alpha1.WorkloadVendor
	pos := strings.LastIndex(vendorName, "/")
	path := fmt.Sprintf("%s/%s/workloadVendor/%s.yaml", vendorDir, vendorName[:pos+1], vendorName[pos+1:])
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("workload.vendor: %s 不存在\n", vendorName))
		return v, err
	}
	err = yaml.Unmarshal(content, &v)
	cuefile := v.Spec
	//替换import为真实内容
	re, _ := regexp.Compile("import\\s*\"([^\"]*)\"")
	matchResult := re.FindAllStringSubmatch(cuefile, -1)
	for _, vv := range matchResult {
		if len(matchResult) > 0 {
			includeMod, err := modTemplate(vendorName, vv[1], vendorDir)
			if err != nil {
				klog.Errorln(err.Error())
				return v, err
			}
			cuefile = strings.ReplaceAll(cuefile, vv[0], includeMod)
		}
	}
	v.Spec = cuefile
	return v, err
}

//获取trait
func GetTrait(name, vendorDir string) (v1alpha1.Trait, error) {
	var err error
	var t v1alpha1.Trait
	pos := strings.LastIndex(name, "/")
	path := fmt.Sprintf("%s/%s/trait/%s.yaml", vendorDir, name[:pos+1], name[pos+1:])
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("trait: %s 不存在\n", name))
		return t, err
	}
	//解析为结构体
	err = yaml.Unmarshal(content, &t)
	return t, err
}
func GetProperties(properties map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{}, 0)
	for k, v := range properties {
		ret[k] = GetValue(v)
	}
	return ret
}

//校验type参数
func CheckTypeParam (workload v1alpha1.Workload, vendorDir string) error{
	var t v1alpha1.WorkloadType
	var err error
	properties := GetProperties(workload.Properties)
	if workload.Type == "" {
		err = errors.New("workload.Type 不能为空")
		return err
	}
	if workload.Vendor == "" {
		err = errors.New("workload.Vendor 不能为空")
		return err
	}
	t, err = GetWorkloadType(workload.Type, vendorDir)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	t, err = GetWorkloadType(workload.Type, vendorDir)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	properties2, err := json.Marshal(properties)
	if err != nil {
		return err
	}
	parameterStr := fmt.Sprintf("parameter:{ \n%s\n}\nparameter:{\n%s\n}", t.Spec.Parameter, string(properties2))
	fileName := RandomString(parameterStr)
	path := fmt.Sprintf("/tmp/%s.cue", fileName)
	ioutil.WriteFile(path, []byte(parameterStr), 0644)
	command := fmt.Sprintf("/usr/local/bin/cue vet -c %s", path)
	cmd := exec.Command("bash", "-c", command)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		klog.Errorln("type参数校验失败: " + stderr.String())
		return errors.New("type参数校验失败: " + stderr.String())
	}
	return nil
}
//校验trait参数
func CheckTraitParam (workloadTrait Trait, vendorDir string) error {
	properties := GetProperties(workloadTrait.Properties)
	properties2, err := json.Marshal(properties)
	if err != nil {
		klog.Errorln(err)
		return errors.New("trait参数序列化失败")
	}
	file, err := GetTrait(workloadTrait.Type, vendorDir)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	tmpcue := fmt.Sprintf("parameter: \n%s\nparameter: {\n%s\n}", string(properties2), file.Spec.Parameter)
	path := fmt.Sprintf("/tmp/%s.cue", RandomString(tmpcue))
	err = ioutil.WriteFile(path, []byte(tmpcue), 0644)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	command := fmt.Sprintf("/usr/local/bin/cue vet -c %s", path)
	cmd := exec.Command("bash", "-c", command)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		klog.Errorln("trait参数校验失败: " + stderr.String())
		return errors.New("trait参数校验失败: " + stderr.String())
	}
	return nil
}

//解析数据
func GetValue(v interface{}) interface{} {
	vType := reflect.TypeOf(v)
	if vType.Kind() == reflect.String {
		after := v.(string)
		return after
	} else if vType.Kind() == reflect.Int {
		after := v.(int)
		return after
	} else if vType.Kind() == reflect.Slice {
		var after []interface{}
		for _, item := range v.([]interface{}) {
			itemValue := GetValue(item)
			after = append(after, itemValue)
		}
		return after
	} else if vType.Kind() == reflect.Struct {
		//todo
		var after interface{}
		return after
	} else if vType.Kind() == reflect.Map {
		after := make(map[string]interface{}, 0)
		for key, val := range v.(map[interface{}]interface{}) {
			newKey := fmt.Sprintf("%s", key)
			newValue := GetValue(val)
			after[newKey] = newValue
		}
		return after
	}
	//todo
	klog.Errorln("其他类型")
	return nil
}
