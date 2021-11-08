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
	Dependencies Dependency `json:"Dependencies"`
	RootDomain   string       `json:"RootDomain"`
	WorkloadPath string       `json:"WorkloadPath"`
}

func PostManifestHandlerFunc(c *gin.Context) {
	var err error
	p := Params{}
	err = c.BindJSON(&p)
	if err != nil {
		fmt.Println(err.Error())
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
	//err = ioutil.WriteFile("tmp/vela.json", str, 0644)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

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
//应用之间的依赖
type ApplicationDependency struct {
	Authorization []Authorization `json:"authorization"`
	ServiceEntry []ServiceEntry `json:"serviceEntry"`
}

//由manifest.yaml生成vale.yaml
func GenValeYaml(instanceId string, application v1alpha1.Application, userconfig string, rootDomain string, dependencies Dependency) (VelaYaml, error) {
	var vela = VelaYaml{"", make(map[string]interface{}, 0)}
	var err error
	vela.Name = application.Metadata.Name

	//traits:ingress的组件
	serviceEntryName := entryService(application.Spec.Workloads)
	authorizationData, serviceEntryData, configmapData, err := parseDependencies(application, dependencies)
	if err != nil {
		return vela, err
	}
	var applicationDependency ApplicationDependency
	applicationDependency.Authorization = authorizationData
	applicationDependency.ServiceEntry = serviceEntryData

	//应用内部的授权
	authorization := make([]Authorization, 0)
	//为每个 service 创建一个 authorization，授权当前应用下的其他服务有访问的权限
	for _, workload := range application.Spec.Workloads {
		authorization = append(authorization,
			Authorization{
				Namespace: instanceId,
				Service:   workload.Name,
				Resources: make([]DependencyUseItem, 0),
			},
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
		service := serviceVela(workload, instanceId, authorization, applicationDependency, configItemData, rootDomain, serviceEntryName)
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
		//fmt.Println(workloadParam[k])
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
	serviceData, err := json.Marshal(workload)
	if err != nil {
		klog.Errorln("vela.Services 序列化失败")
		return cmdResult, errors.New("vela.Services 序列化失败")
	}
	content := fmt.Sprintf(finnnalCueFileContent, ctxObjData, serviceData, template)
	fileName := RandomString(content)
	path := fmt.Sprintf("/tmp/%s.cue", fileName)
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
	path := fmt.Sprintf("%s/%s/%s.cue", vendorDir, workloadVendor[:pos+1], mod)
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
func serviceVela(workload v1alpha1.Workload, instanceid string, authorization []Authorization, ApplicationDependency ApplicationDependency, configItemData []ConfigItemDataItem, rootDomain string, serviceEntryName string) interface{} {
	properties := GetProperties(workload.Properties)
	properties["authorization"] = authorization
	configs2 := make([]interface{}, 0)
	if configs, ok := properties["configs"]; ok {
		for _, v := range configs.([]interface{}) {
			configs2 = append(configs2, v)
		}
	}
	configs2 = append(configs2, ConfigItem{"/etc/configs", "", configItemData})
	properties["configs"] = configs2

	//整合trait参数
	type Trait struct {
		Type       string	`json:"type"`
		Properties v1alpha1.Properties `json:"properties"`
	}
	var traits = make(map[string]interface{}, 0)
	if len(workload.Traits) > 0 {
		for _, trait:= range workload.Traits {
			if trait.Type == "globalsphare.com/v1alpha1/trait/ingress" {
				traitProperties := make(map[string]interface{}, 0)
				traitProperties["host"] = fmt.Sprintf("%s.%s", instanceid, rootDomain)
				path := make([]string, 0)
				traitProperties["path"] = append(path, "/*")
				traits[trait.Type] = traitProperties
				//添加一个外部依赖的trait
				traits["globalsphare.com/v1alpha1/trait/dependency"] = ApplicationDependency
			}else{
				traits[trait.Type] =  GetProperties(trait.Properties)
			}
		}
		properties["traits"] = traits
	}
	return properties
}

//处理依赖
func parseDependencies(application v1alpha1.Application, dependencies Dependency) ([]Authorization, []ServiceEntry, map[string]string, error) {
	var err error
	auth := make([]Authorization, 0)
	//外部服务调用
	svcEntry := make([]ServiceEntry, 0)
	//运行时配置
	cm := make(map[string]string, 0)
	allDependency := make(map[string][]DependencyUseItem)
	for _,j := range application.Spec.Dependencies{
		resource, err := ApiParse(j.Items) //[]DependencyUseItem
		if err != nil {
			klog.Errorln(err)
			return auth, svcEntry, cm, err
		}
		allDependency[j.Name] = resource
	}

	//从manifest.yaml中解析uses
	for _,v := range dependencies.Internal {
		auth = append(auth, Authorization{
			v.Instanceid, v.EntryService, allDependency[v.Name],
		})
		cm[v.Name] = fmt.Sprintf("%s.%s.svc.cluster.local", v.EntryService, v.Instanceid)
	}
	for _,item := range dependencies.External {
		arr, err := url.ParseRequestURI(item.Location)
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
	return auth, svcEntry, cm, err
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
		//trait:ingress只能有一个
		//检查trait参数
		traitCount := 0
		if len(workload.Traits) > 0 {
			for _, trait := range workload.Traits {
				err = CheckTraitParam(trait, vendorDir)
				if err != nil {
					return returnData, err
				}
				arr := strings.Split(trait.Type, "/")
				if arr[len(arr)-1] == "ingress" {
					traitCount++
				}
			}
		}
		if traitCount > 1 {
			err = errors.New("检测到多个ingress")
			return returnData, err
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
	path := fmt.Sprintf("%s/%s.yaml", vendorDir, typeName)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("workload.Type: %s 不存在\n", path))
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
	path := fmt.Sprintf("%s/%s.yaml", vendorDir, vendorName)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("workload.vendor: %s 不存在\n", path))
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
	path := fmt.Sprintf("%s/%s.yaml", vendorDir, name)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("trait: %s 不存在\n", path))
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
