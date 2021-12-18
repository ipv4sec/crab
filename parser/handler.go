package parser

import (
	"crab/aam/v1alpha1"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"strings"
)

type Params struct {
	Content      string      `json:"Content"`
	Instanceid   string      `json:"InstanceId"`
	Userconfig   interface{} `json:"UserConfig"`
	Dependencies Dependency  `json:"Dependencies"`
	Host         string      `json:"Host"`
	WorkloadPath string      `json:"WorkloadPath"`
}
type Result struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

func PostManifestHandlerFunc(c *gin.Context) {
	var err error
	p := Params{}
	err = c.BindJSON(&p)
	if err != nil {
		klog.Infoln(err)
		c.JSON(200, Result{ErrBadRequest, "参数格式错误"})
		return
	}
	if p.Content == "" || p.Instanceid == "" {
		c.JSON(200, Result{ErrBadRequest, "缺少参数"})
		return
	}
	userconfig, err := json.Marshal(p.Userconfig)
	if err != nil {
		c.JSON(200, Result{ErrInternalServer, "运行时配置序列化失败"})
		return
	}
	userconfigStr := strings.TrimSpace(string(userconfig))
	if userconfigStr == "null" || userconfigStr == "" {
		userconfigStr = "{}"
	}
	//解析描述文件
	var application v1alpha1.Application
	err = yaml.Unmarshal([]byte(p.Content), &application)
	if err != nil {
		c.JSON(200, Result{ErrBadRequest, "描述文件解析失败"})
		return
	}

	//验证参数，返回参数json,返回vendor内容
	workloadResource, err := checkParams(application, p.WorkloadPath)
	if err != nil {
		c.JSON(200, Result{ErrBadRequest, err.Error()})
		return
	}

	//生成vale.yaml文件
	vale, err := GenValeYaml(p.Instanceid, application, userconfigStr, p.Host, p.Dependencies)
	if err != nil {
		c.JSON(200, Result{ErrInternalServer, err.Error()})
		return
	}
	str, err := json.Marshal(vale)
	if err != nil {
		klog.Errorln(err)
		return
	}
	tmpName := fmt.Sprintf("/tmp/%s-vela.json", RandomStr())
	ioutil.WriteFile(tmpName, str, 0644)

	//生成k8s.yaml文件
	k8s, err := GenK8sYaml(p.Instanceid, vale, workloadResource)
	if err != nil {
		klog.Errorln(err)
		c.JSON(200, Result{ErrInternalServer, err.Error()})
		return
	}
	tmpName = fmt.Sprintf("/tmp/%s-k8s.yaml", RandomStr())
	ioutil.WriteFile(tmpName, []byte(k8s), 0644)
	c.JSON(200, Result{0, k8s})
}

func GetSystemTemplateFunc(c *gin.Context) {
	text, err := ioutil.ReadFile("assets/cue/systemTemplate.cue")
	if err != nil {
		fmt.Println(err)
		c.JSON(200, Result{
			Code:   ErrInternalServer,
			Result: "模板不存在",
		})
		return
	}
	c.JSON(200, Result{
		Code:   0,
		Result: string(text),
	})
}
