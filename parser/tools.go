package parser

import (
	"crab/aam/v1alpha1"
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gojek/heimdall/v7/httpclient"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var HTTPClient = httpclient.NewClient(httpclient.WithHTTPTimeout(time.Second * 30))

//由manifest.yaml生成vale.yaml
func GenValeYaml(instanceId string, application v1alpha1.Application, userconfigs string, host string, dependencies Dependency) (VelaYaml, error) {
	var vela = VelaYaml{"", make(map[string]interface{}, 0)}
	var err error
	vela.Name = application.Metadata.Name

	authorization, serviceEntry, err := parseDependencies(application, dependencies)
	if err != nil {
		return vela, err
	}
	//应用内部的授权
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
	//依赖内部应用的host
	dependHost := make(dependencyHost, 0)
	for _, v := range dependencies.Internal {
		//host
		dependHost[v.Name] = dependencyHostItem{
			fmt.Sprintf("%s.%s.svc.cluster.local", v.EntryService, v.Instanceid),
		}
		//授权
		authorization = append(authorization,
			Authorization{
				Namespace: v.Instanceid,
				Service:   v.EntryService,
				Resources: make([]DependencyUseItem, 0),
			},
		)
	}
	for _, workload := range application.Spec.Workloads {
		properties := GetProperties(workload.Properties)
		properties["authorization"] = authorization
		properties["serviceEntry"] = serviceEntry
		properties["userconfigs"] = userconfigs
		properties["dependencies"] = dependHost
		//整合trait参数
		traitList := make([]string, 0)
		if len(workload.Traits) > 0 {
			for _, trait := range workload.Traits {
				arr := strings.Split(trait.Type, "/")
				shortName := arr[len(arr)-1]
				traitList = append(traitList, shortName)
				if shortName == "ingress" {
					ingressProperties := make(map[string]interface{}, 0)
					if host == ""{
						return vela, errors.New("访问域名不能为空")
					}
					ingressProperties["host"] = host
					ingressProperties["path"] = []string{"/*"}
					properties[shortName] = ingressProperties
				} else {
					properties[shortName] = GetProperties(trait.Properties)
				}
			}
		}
		properties["traits"] = traitList
		vela.Services[workload.Name] = properties
	}
	return vela, nil
}

//由vale.yaml生成k8s
func GenK8sYaml(instanceId string, vela VelaYaml, workloadParam map[string]WorkloadParam) (string, error) {
	//先创建命名空间
	finalContext := `
apiVersion: v1
kind: Namespace
metadata:
  labels:
    istio-injection: enabled
  name: %s
---
`
	finalContext = fmt.Sprintf(finalContext, instanceId)
	//自动追加的部分
	//处理workload
	for k, v := range vela.Services {
		ctx := ContextObj{
			vela.Name,
			k,
			instanceId,
		}
		k8sStr, err := Export(ctx, workloadParam[k], v)
		if err != nil {
			klog.Errorln(err)
			return "", err
		}

		//if strings.Index(k8sStr, "---") != -1 {
		//	s := strings.Split(k8sStr, "---")
		//	for _, v := range s {
		//		v = strings.TrimSpace(v)
		//		if v == "" {
		//			continue
		//		}
		//		if strings.Index(finalContext, v) == -1 {
		//			finalContext += v + "\n---\n"
		//		}
		//	}
		//} else {
		finalContext += k8sStr + "\n---\n"
		//}
	}
	finalContext = strings.Trim(strings.TrimSpace(finalContext), "---")
	finalContext = fmt.Sprintf("# appName: %s\n%s", vela.Name, finalContext)
	return finalContext, nil
}

func Export(ctxObj ContextObj, workloadParam WorkloadParam, workload interface{}) (string, error) {
	var k8s = ""
	template := workloadParam.VendorCue
	ctxData := make(map[string]interface{}, 0)
	ctxData["context"] = ctxObj
	ctxObjData, err := json.Marshal(ctxData)
	if err != nil {
		klog.Errorln("ctxObj 序列化失败: ", err.Error())
		return "", errors.New("ctxObj 序列化失败")
	}
	serviceData, err := json.Marshal(workload)
	if err != nil {
		klog.Errorln("vela.Services 序列化失败: ", err.Error())
		return "", errors.New("vela.Services 序列化失败")
	}
	cueStr := fmt.Sprintf("%s\nparameter:%s\n%s", ctxObjData, serviceData, template)
	err = ioutil.WriteFile(fmt.Sprintf("/tmp/%s-%s.cue", ctxObj.Namespace, ctxObj.WorkloadName), []byte(cueStr), 0644)
	if err != nil {
		klog.Errorln("保存cue文件错误: ", err.Error())
		return "", err
	}
	//处理cue内置的pkg
	cueStr = MoveCuePkgToTop(cueStr)
	var ctx *cue.Context
	var value cue.Value
	ctx = cuecontext.New()
	value = ctx.CompileString(cueStr)
	if value.Err() != nil {
		klog.Errorln("cue生成yaml失败: ", value.Err().Error())
		return "", value.Err()
	}
	context := make(map[string]interface{}, 0)
	err = value.Decode(&context)
	for k,v := range context {
		if k != "context" && k != "parameter" && k != "namespace" {
			b, err := yaml.Marshal(v)
			if err != nil {
				klog.Errorln("解析context失败: ", err)
				return "", err
			}
			k8s = fmt.Sprintf("%s\n---\n%s", k8s, string(b))
		}
	}
	return strings.TrimSpace(k8s), nil
}

//处理依赖
func parseDependencies(application v1alpha1.Application, dependencies Dependency) ([]Authorization, []ServiceEntry, error) {
	var err error
	auth := make([]Authorization, 0)
	//外部服务调用
	svcEntry := make([]ServiceEntry, 0)
	for _, item := range dependencies.External {
		var host, address string
		arr, err := url.ParseRequestURI(item.Location)
		if err != nil {
			klog.Errorln(err.Error())
			return auth, svcEntry, err
		}
		var protocol string
		if arr.Scheme == "https" {
			protocol = "TLS"
		} else if arr.Scheme == "http" {
			protocol = "http"
		} else {
			klog.Errorln("protocol of the location is not http or https.")
			return auth, svcEntry, errors.New("protocol of the location is not http or https.")
		}
		hostArr := strings.Split(arr.Host, ":")
		var port int
		if len(hostArr) == 1 { //没有指定端口号
			if protocol == "http" {
				port = 80
			} else {
				port = 443
			}
		} else { //指定端口号
			port, err = strconv.Atoi(hostArr[1])
			if err != nil {
				klog.Errorln("转int失败")
				return auth, svcEntry, errors.New("转int失败")
			}
		}
		ipAddress := net.ParseIP(hostArr[0])
		if ipAddress != nil { //ip
			host = fmt.Sprintf("serviceEntry-%s-%s", application.Metadata.Name, item.Name)
			address = ipAddress.String()
		} else {
			host = arr.Host
		}
		svcEntry = append(svcEntry, ServiceEntry{item.Name, address, host, port, protocol})
	}
	return auth, svcEntry, err
}

//获取WorkloadType
func GetWorkloadType(typeName string) (v1alpha1.WorkloadType, error) {
	var t v1alpha1.WorkloadType
	value, err := GetWorkloadDef("workloadType", typeName)
	if err != nil {
		klog.Errorln("获取type失败: ", err.Error())
		return t, err
	}
	err = yaml.Unmarshal([]byte(value), &t)
	return t, err
}

//获取trait
func GetTrait(name string) (v1alpha1.Trait, error) {
	var t v1alpha1.Trait
	value, err := GetWorkloadDef("trait", name)
	if err != nil {
		return t, err
	}
	//解析为结构体
	err = yaml.Unmarshal([]byte(value), &t)
	return t, err
}

//获取WorkloadVendor
func GetWorkloadVendor(name string) (v1alpha1.WorkloadVendor, error) {
	var v v1alpha1.WorkloadVendor
	value, err := GetWorkloadDef("workloadVendor", name)
	if err != nil {
		return v, err
	}
	err = yaml.Unmarshal([]byte(value), &v)
	return v, err
	//path := fmt.Sprintf("%s%s.yaml", workloadPath, vendorName)
	//content, err := ioutil.ReadFile(path)
	//if err != nil {
	//	err = errors.New(fmt.Sprintf("workload.vendor: %s 不存在\n", path))
	//	return v, err
	//}
	//err = yaml.Unmarshal(content, &v)
	//cuefile := v.Spec
	//替换import为真实内容
	//re, _ := regexp.Compile("import\\s*\"([^\"]*)\"")
	//matchResult := re.FindAllStringSubmatch(cuefile, -1)
	//for _, vv := range matchResult {
	//	if len(matchResult) > 0 {
	//		if _, ok := cuePkg[vv[1]]; ok {
	//			continue
	//		}
	//		includeMod, err := modTemplate(workloadPath, vendorName, vv[1])
	//		if err != nil {
	//			klog.Errorln(err.Error())
	//			return v, err
	//		}
	//		cuefile = strings.ReplaceAll(cuefile, vv[0], includeMod)
	//	}
	//}
	//v.Spec = cuefile
}

//获取cue模板
func modTemplate(workloadPath, vendorDir, mod string) (string, error) {
	var err error
	pos := strings.LastIndex(vendorDir, "/")
	path := fmt.Sprintf("%s%s%s.cue", workloadPath, vendorDir[:pos+1], mod)
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
			if _, ok := cuePkg[v[1]]; ok {
				continue
			}
			includeMod, err := modTemplate(workloadPath, vendorDir, v[1])
			if err != nil {
				klog.Errorln(err.Error())
				return "", err
			}
			content = strings.ReplaceAll(content, v[0], includeMod)
		}
	}
	return content, nil
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

//cue内置的pkg，放到cue文件第一行
func MoveCuePkgToTop(str string) string {
	pkg := make([]string, 0)
	re, _ := regexp.Compile("import\\s*\"([^\"]*)\"")
	matchResult := re.FindAllStringSubmatch(str, -1)
	for _, v := range matchResult {
		if len(matchResult) > 0 {
			if _, ok := cuePkg[v[1]]; ok {
				pkg = append(pkg, v[0])
				str = strings.ReplaceAll(str, v[0], "")
			}
		}
	}
	return strings.Join(pkg, "\n") + "\n" + str
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

//校验trait参数
func CheckTraitParam(workloadTrait Trait) error {
	properties := GetProperties(workloadTrait.Properties)
	properties2, err := json.Marshal(properties)
	if err != nil {
		klog.Errorln(err)
		return errors.New("trait参数序列化失败")
	}
	file, err := GetTrait(workloadTrait.Type)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	cueStr := fmt.Sprintf("parameter:%s\nparameter: {\n%s\n}", string(properties2), file.Spec.Parameter)
	var ctx *cue.Context
	var value cue.Value
	ctx = cuecontext.New()
	value = ctx.CompileString(cueStr)
	err = value.Validate(cue.Concrete(true))
	if err != nil {
		return err
	}
	return nil
}

//校验type参数
func CheckTypeParam(workload v1alpha1.Workload) error {
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
	t, err = GetWorkloadType(workload.Type)
	if err != nil {
		klog.Infoln(err)
		return err
	}
	properties2, err := json.Marshal(properties)
	if err != nil {
		return err
	}
	parameterStr := fmt.Sprintf("parameter:{ \n%s\n}\nparameter:{\n%s\n}", t.Spec.Parameter, string(properties2))
	var ctx *cue.Context
	var value cue.Value
	ctx = cuecontext.New()
	value = ctx.CompileString(parameterStr)
	err = value.Validate(cue.Concrete(true))
	if err != nil {
		klog.Errorln(err)
		return err
	}
	return nil
}

func checkParams(application v1alpha1.Application) (map[string]WorkloadParam, error) {
	var err error
	returnData := make(map[string]WorkloadParam, 0)
	if len(application.Spec.Workloads) == 0 {
		err = errors.New("application.Spec.Workloads 不能为空")
		return returnData, err
	}
	ingressCount := 0
	for _, workload := range application.Spec.Workloads {
		//检查type参数
		err = CheckTypeParam(workload)
		if err != nil {
			klog.Errorln(err)
			return returnData, err
		}
		//检查trait参数
		if len(workload.Traits) > 0 {
			for _, trait := range workload.Traits {
				err = CheckTraitParam(trait)
				if err != nil {
					klog.Errorln(err)
					return returnData, err
				}
				arr := strings.Split(trait.Type, "/")
				if arr[len(arr)-1] == "ingress" {
					ingressCount++
				}
			}
		}
		var workloadParams WorkloadParam
		workloadParams.Type = workload.Type
		workloadParams.Vendor = workload.Vendor

		properties := GetProperties(workload.Properties)
		workloadParams.Parameter = properties

		t, _ := GetWorkloadType(workload.Type)
		workloadParams.Traits = t.Spec.Traits

		var v v1alpha1.WorkloadVendor
		v, err = GetWorkloadVendor(workload.Vendor)
		if err != nil {
			klog.Errorln(err)
			return returnData, err
		}
		workloadParams.VendorCue = v.Spec
		returnData[workload.Name] = workloadParams
	}
	//trait:ingress最多一个
	if ingressCount > 1 {
		err = errors.New("检测到多个ingress")
		return returnData, err
	}
	return returnData, nil
}
func GetProperties(properties map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{}, 0)
	for k, v := range properties {
		ret[k] = GetValue(v)
	}
	return ret
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//获取workload定义
func GetWorkloadDef(kind, name string) (string, error) {
	type def struct {
		Id         int `json:"id"`
		Name       string `json:"name"`
		ApiVersion string `json:"apiVersion"`
		Value      string `json:"value"`
		Type       int `json:"type"`
	}
	var err error
	kind = strings.TrimSpace(kind)
	if kind == "" {
		return "", errors.New("kind不能为空")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("name不能为空")
	}
	res, err := HTTPClient.Get(fmt.Sprintf("http://127.0.0.1:3000/%s/%s", kind, name), nil)
	if err != nil {
		return "", fmt.Errorf("请求api错误: %w", err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取api返回错误: %w", err)
	}
	klog.Infoln("api返回内容:", string(bodyBytes))
	var reply struct {
		Code   int         `json:"code"`
		Result def `json:"result"`
	}
	err = json.Unmarshal(bodyBytes, &reply)
	if err != nil {
		return "", fmt.Errorf("api返回序列化错误: %w", err)
	}
	if reply.Code != 0 {
		return "", fmt.Errorf("api返回错误: %v", reply.Result)
	}
	return reply.Result.Value, err
}