package parser

import (
	"crab/aam/v1alpha1"
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
	"regexp"
	"strings"
	"strconv"
)

type contextObj struct {
	AppName       string `json:"appName"`
	ComponentName string `json:"componentName"`
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
	Dependencies []v1alpha1.Dependency `json:"Dependencies"`
	RootDomain string `json:"RootDomain"`
	WorkloadPath string `json:"WorkloadPath"`
}

//验证type,vendor返回的数据
type WorkloadParams struct{
	Parameter string `json:"parameter"`
	Type string `json:"type"`
	Vendor string `json:"vendor"`
	VendorCue string `json:"vendorCue"`
	Traits []string `json:"traits"`
}

func PostManifestHandlerFunc(c *gin.Context) {
	var err error
	p := Params{}
	err = c.BindJSON(&p)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数格式错误"))
		return
	}
	if p.Content == "" || p.Instanceid == "" || p.RootDomain == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	userconfigStr, err := json.Marshal(p.Userconfig)
	if err != nil {
		c.JSON(200, result{10101,"序列化失败"})
		return
	}
	//解析描述文件
	var application v1alpha1.Application
	err = yaml.Unmarshal([]byte(p.Content), &application)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "文件解析失败"))
		return
	}
	fmt.Println("---application---")
	fmt.Printf("%+v\n", application)

	//for i := 0; i < len(application.Spec.Workloads); i++ {
	//	for k, v := range application.Spec.Workloads[i].Properties {
	//		klog.Infoln(reflect.ValueOf(k))
	//		klog.Infoln(reflect.ValueOf(v))
	//	}
	//}


	//验证参数，返回参数json,返回vendor内容
	//test
	if p.WorkloadPath == "" {
		p.WorkloadPath = "/Users/huanqiu/Desktop/uploads"
	}


	workloadResource, err := checkParams(application, p.WorkloadPath)
	//_ = workloadResource
	//if err != nil {
	//	c.JSON(200, err.Error())
	//	return
	//}

	//生成vale.yaml文件
	vale, err := GenValeYaml(p.Instanceid, application, string(userconfigStr), p.RootDomain, p.Dependencies)
	if err != nil {
		returnData := result{
			10101,
			err.Error(),
		}
		c.JSON(200, returnData)
		return
	}
	str,err := json.Marshal(vale)
	if err != nil {
		klog.Errorln(err)
		return
	}
	ioutil.WriteFile("/tmp/vela.json", str, 0644)

	//生成k8s.yaml文件
	k8s, err := GenK8sYaml(p.Instanceid, vale, workloadResource)
	if err != nil {
		returnData := result{
			10101,
			err.Error(),
		}
		klog.Errorln(err.Error())
		c.JSON(200, returnData)
		return
	}
	k8s2,err := yaml.Marshal(k8s)
	if err != nil {
		fmt.Println(err)
		return
	}
	ioutil.WriteFile("tmp/k8s.yaml", k8s2, 0644)
	returnData := struct {
		Code   int    `json:"code"`
		Result string `json:"result"`
	}{
		0,
		"ok",
	}
	c.JSON(200, returnData)
}
//由manifest.yaml生成vale.yaml
func GenValeYaml(instanceId string, application v1alpha1.Application, userconfig string, rootDomain string, dependencies []v1alpha1.Dependency) (VelaYaml, error) {
	var vela = VelaYaml{"", make(map[string]interface{}, 0)}
	var err error
	vela.Name = application.Metadata.Name
	fmt.Println("应用名称：", vela.Name)

	//Workloads
	if len(application.Spec.Workloads) == 0 {
		return vela, errors.New("workloads不能为空")
	}

	//traits:ingress的组件
	serviceEntryName := entryService(application.Spec.Workloads)
	//fmt.Println("ingress：", serviceEntryName)
	authorizationData, serviceEntryData, configmapData, err := parseDependencies(dependencies)
	if err != nil {
		return vela, err
	}

	//为每个 service 创建一个 authorization，授权当前应用下的其他服务有访问的权限
	for _, component := range application.Spec.Workloads {
		authorizationData = append(authorizationData,
			v1alpha1.Authorization{
				Namespace: instanceId,
				Service:   component.Name,
				Resources: make([]v1alpha1.DependencyUseItem, 0)},
		)
	}

	//configmap
	configItemData := make([]v1alpha1.ConfigItemDataItem, 0)
	for k, v := range configmapData {
		configItemData = append(configItemData, v1alpha1.ConfigItemDataItem{Name: fmt.Sprintf("%s.host", k), Value: v})
	}
	//添加应用时填写的运行时配置
	configItemData = append(configItemData, v1alpha1.ConfigItemDataItem{Name: "userconfig", Value: userconfig})
	for _, svc := range application.Spec.Workloads {
		service := serviceVela(svc, instanceId, authorizationData, serviceEntryData, configItemData, rootDomain, serviceEntryName)
		vela.Services[svc.Name] = service
	}

	return vela, nil
}

//由vale.yaml生成k8s
func GenK8sYaml(instanceid string, vela VelaYaml, workloadParams map[string]WorkloadParams) (ParserData, error) {
	parserData := ParserData{
		Init:      "",
		Name:      "",
		Workloads:make(map[string]Workload,0),
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
	ns = fmt.Sprintf(ns, instanceid, vela.Name, instanceid,instanceid,instanceid)
	parserData.Init = ns
	parserData.Name = vela.Name
	//处理workload
	for k, v := range vela.Services {
		ctxObj := make(map[string]contextObj, 0)
		ctxObj["context"] = contextObj{
			vela.Name,
			k,
			instanceid,
		}
		finnnalCueFileContent := "%s\nparameter:%s\n%s"
		template := workloadParams[k].VendorCue
		ctxObjData, err := json.Marshal(ctxObj)
		if err != nil {
			klog.Errorln("ctxObj 序列化失败")
			return parserData, errors.New("ctxObj 序列化失败")
		}
		serviceItem, err := json.Marshal(v)
		if err != nil {
			klog.Errorln("vela.Services 序列化失败")
			return parserData, errors.New("vela.Services 序列化失败")
		}
		content := fmt.Sprintf(finnnalCueFileContent, ctxObjData, serviceItem, template)
		fileName := RandomString(content)
		path := fmt.Sprintf("/tmp/%s.cue", fileName)
		err = ioutil.WriteFile(path, []byte(content), 0644)
		if err != nil {
			klog.Errorln(err.Error())
			return parserData, err
		}
		command := fmt.Sprintf("/usr/local/bin/cue export -f %s", path)
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			klog.Errorln("执行命令错误", err.Error())
			return parserData, err
		}
		var cmdResult map[string]map[string]interface{}
		err = json.Unmarshal(output, &cmdResult)
		if err != nil {
			klog.Errorln(err.Error())
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
			err = errors.New("vendor中无construct")
			return parserData, err
		}
		workload.Construct = construct
		traits := make(map[string]string, 0)
		if len(workloadParams[k].Traits) >0 {//有trait
			for _,v := range workloadParams[k].Traits {
				count = 0
				arr := strings.Split(v,"/")
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
				if count == 0{
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

//获取cue模板
func modTemplate(workloadVendor, mod,vendorDir string) (string, error) {
	var err error
	pos := strings.LastIndex(workloadVendor, "/")
	templatePath := vendorDir +"/" + workloadVendor[:pos+1]+mod+".cue"
	if !FileExist(templatePath) {
		return "", errors.New(fmt.Sprintf("文件：%s 不存在", templatePath))
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
	h.Write([]byte(str + string(rand.Intn(1000))))
	return hex.EncodeToString(h.Sum(nil))
}

//生成kubevela格式的service
func serviceVela(svc v1alpha1.Workload, instanceid string, authorization []v1alpha1.Authorization, serviceentry []v1alpha1.ServiceEntry, configItemData []v1alpha1.ConfigItemDataItem, rootDomain string, serviceEntryName string) interface{} {
	fmt.Println(svc.Type)
	arr := strings.Split(svc.Type, "/")
	workloadType := arr[len(arr)-1]
	if workloadType == "webservice" {
		service := WebserviceVela{
			Workload:      svc.Type,
			Type:          svc.Type,
			Image:         svc.Properties.Image,
			Configs: 	   append(svc.Properties.Configs, v1alpha1.ConfigItem{"/etc/configs", "", configItemData}),
			Storage:       svc.Properties.Storage,
			Init:          svc.Properties.Init,
			After:         svc.Properties.After,
			Port:          0,
			Cmd:           svc.Properties.Cmd,
			Args:          svc.Properties.Args,
			Env:           make([]v1alpha1.EnvItem, 0),
			Traits:        svc.Traits,
			Authorization: authorization,
			Serviceentry:  serviceentry,
			Namespace:     instanceid,
			Entry:         v1alpha1.Entry{},
		}
		if serviceEntryName == svc.Name {
			path := make([]string, 0)
			path = append(path, "/*")
			service.Entry = v1alpha1.Entry{
				fmt.Sprintf("%s.%s", instanceid, rootDomain),
				path,
			}
		} else {
			service.Entry = v1alpha1.Entry{
				"",
				make([]string, 0),
			}
		}
		return service
	} else if workloadType == "worker" {
		service := WorkerVela{
			Workload:      svc.Type,
			Type:          svc.Type,
			Image:         svc.Properties.Image,
			Cmd:           svc.Properties.Cmd,
			Args:          svc.Properties.Args,
			Env:           make([]v1alpha1.EnvItem, 0),
			After:         svc.Properties.After,
			Init:          svc.Properties.Init,
			Configs:       append(svc.Properties.Configs, v1alpha1.ConfigItem{"/etc/configs", "", configItemData}),
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
	} else if workloadType == "mysql" {
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
	} else if workloadType == "redis" {
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
func parseDependencies(dependencies []v1alpha1.Dependency) ([]v1alpha1.Authorization, []v1alpha1.ServiceEntry, map[string]string, error) {
	var err error
	authorization := make([]v1alpha1.Authorization, 0)
	serviceEntry := make([]v1alpha1.ServiceEntry, 0)
	configmap := make(map[string]string, 0)
	//解析uses
	dependencyVelas := make([]v1alpha1.DependencyVela, 0)
	for _, v := range dependencies {
		if v.Instanceid != "" && v.EntryService == "" {
			return authorization, serviceEntry, configmap, errors.New("dependencies.entryService不能为空")
		}
		resource,err := ApiParse(v.Items)
		if err != nil {
			return authorization, serviceEntry, configmap, err
		}
		dependencyVelas = append(dependencyVelas, v1alpha1.DependencyVela{
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
func dependendService(dependencyVelas []v1alpha1.DependencyVela) ([]v1alpha1.Authorization, []v1alpha1.ServiceEntry, map[string]string, error) {
	auth := make([]v1alpha1.Authorization, 0)
	//外部服务调用
	svcEntry := make([]v1alpha1.ServiceEntry, 0)
	//运行时配置
	cm := make(map[string]string, 0)

	for _, v := range dependencyVelas {
		if v.Instanceid != "" {
			auth = append(auth, v1alpha1.Authorization{
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
				auth = append(auth, v1alpha1.Authorization{arr[0], arr[1], v.Resource})
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
					v1alpha1.ServiceEntry{arr.Host, port, protocol},
				)
			}
		}
	}
	return auth, svcEntry, cm, nil
}

//traits中包含ingress的组件名称
func entryService(workloads []v1alpha1.Workload) string {
	for _, svc := range workloads {
		for _, v := range svc.Traits {
			if v.Type == "ingress" {
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

func ApiParse(uses map[string][]string) ([]v1alpha1.DependencyUseItem, error) {
	var err error
	rtn := make([]v1alpha1.DependencyUseItem, 0)
	for k, v := range uses {
		count := 0
		actions := make([]string, 0)
		for _, option := range v {
			if option == "create" {
				actions = append(actions, "POST")
			}else if option == "read" {
				actions = append(actions, "GET", "HEAD", "OPTIONS")
			}else if option == "update" {
				actions = append(actions, "PUT", "PATCH")
			}else if option == "delete" {
				actions = append(actions, "DELETE")
			}else{
				return rtn, errors.New(fmt.Sprintf("依赖资源的操作类型(%s)不存在\n", option))
			}
			count++
		}
		if count == 0 {
			return rtn, errors.New("依赖资源的操作类型不能为空")
		}
		rtn = append(rtn, v1alpha1.DependencyUseItem{k, actions})
	}
	return rtn,err
}

func checkParams(application v1alpha1.Application, vendorDir string) (map[string]WorkloadParams, error) {
	var err error
	returnData := make(map[string]WorkloadParams, 0)
	if len(application.Spec.Workloads) == 0 {
		err = errors.New("application.Spec.Workloads 不能为空")
		return returnData, err
	}
	for _, workload := range application.Spec.Workloads {
		var workloadParams WorkloadParams
		//fmt.Println("-----workload.Properties-----")
		//fmt.Printf("%+v\n", workload.Properties)


		workloadParams.Traits = make([]string, 0)
		if workload.Type == "" {
			err = errors.New("workload.Type 不能为空")
			return returnData, err
		}
		if workload.Vendor == "" {
			err = errors.New("workload.Vendor 不能为空")
			return returnData,err
		}
		var t v1alpha1.WorkloadType
		t, err = GetWorkloadType(workload.Type,vendorDir)
		if err != nil {
			fmt.Println(err)
			return returnData,err
		}

		workloadParams.Traits = t.Spec.Traits
		workloadParams.Type = workload.Type
		workloadParams.Vendor = workload.Vendor

		var strJson []byte
		var p interface{}
		if t.Metadata.Name == "mysql" {
			p = MysqlParam{
				Init:    workload.Properties.Init,
				Rootpwd: workload.Properties.Rootpwd,
				Storage: struct {
					Capacity string `json:"capacity"`
				}{Capacity:workload.Properties.Storage.Capacity},
				After: workload.Properties.After,
			}
		}else if t.Metadata.Name == "webservice" {
			p = WebserviceParam{
				Image: workload.Properties.Image,
				Port:  workload.Properties.Port,
				Cmd:   workload.Properties.Cmd,
				Args:  workload.Properties.Args,
				Cpu:   workload.Properties.Cpu,
				Env:   workload.Properties.Env,
				After: workload.Properties.After,
			}
		}else if t.Metadata.Name == "worker" {
			p = WorkerParam{
				Image: workload.Properties.Image,
				Port:  workload.Properties.Port,
				Cmd:   workload.Properties.Cmd,
				Args:  workload.Properties.Args,
				Cpu:   workload.Properties.Cpu,
				Env:   workload.Properties.Env,
				After: workload.Properties.After,
			}

		}else if t.Metadata.Name == "redis" {
			p = RedisParam{
				After: workload.Properties.After,
			}
		}
		strJson,err = json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			return returnData, err
		}
		workloadParams.Parameter = string(strJson)

		//检查参数
		parameterStr := fmt.Sprintf("parameter:{ \n%s\n}\nparameter:{\n%s\n}", t.Spec.Parameter, string(strJson))
		fileName := RandomString(parameterStr)
		path := fmt.Sprintf("/tmp/%s.cue", fileName)
		ioutil.WriteFile(path, []byte(parameterStr),0644)
		command := fmt.Sprintf("/usr/local/bin/cue vet %s", path)
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			klog.Errorln("执行命令错误", err.Error())
			return returnData, errors.New("参数错误")
		}
		_ = output

		var v v1alpha1.WorkloadVendor
		v,err = GetWorkloadVendor(workload.Vendor,vendorDir)
		if err != nil {
			return returnData, err
		}
		workloadParams.VendorCue = v.Spec
		returnData[workload.Name] = workloadParams
	}
	return returnData,nil
}

//获取WorkloadType
func GetWorkloadType(typeName,vendorDir string) (v1alpha1.WorkloadType, error){
	var t v1alpha1.WorkloadType
	var err error
	content, err := ioutil.ReadFile(vendorDir + "/" + typeName + ".yaml")
	if err != nil {
		err = errors.New(fmt.Sprintf("workload.Type: %s 不存在\n", typeName))
		return t, err
	}
	//解析为结构体
	err = yaml.Unmarshal(content, &t)
	return t, err
}
//获取WorkloadVendor
func GetWorkloadVendor(vendorName,vendorDir string) (v1alpha1.WorkloadVendor, error){
	var err error
	var v v1alpha1.WorkloadVendor
	content,err := ioutil.ReadFile(vendorDir+"/"+vendorName+".yaml")
	if err != nil {
		fmt.Println(err)
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
			includeMod, err := modTemplate(vendorName, vv[1],vendorDir)
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

func trans(data map[string]interface{}) string {
	var s = ""
	for k,v := range data{
		var key, value string
		switch v.(type) {
		case string:
			key = k
			value = v.(string)
		case interface{}:
			key = k
			var s2 = ""
			for kk,vv := range v.(map[interface{}]interface{}){
				var key, value string
				switch kk.(type) {
				case string:
					key = k
					fmt.Println("lkkkk", key)
				}
				switch vv.(type) {
				case string:
					value = vv.(string)
				}
				s2 = s2 + fmt.Sprintf("\"%s\":\"%s\",", key, value)
			}
			s2 = strings.Trim(s2,",")
			value = "{" + s2 + "}"
		}
		s = s + fmt.Sprintf("\"%s\":\"%s\",", key, value)
	}
	s = strings.Trim(s,",")
	s = "{" + s + "}"
	fmt.Println(s)
	return s
}

