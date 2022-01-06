package tool

import (
	"bytes"
	"crab/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	yaml2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/util/yaml"
	"os/exec"
	"strings"
	"time"
)

type generic struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name      string `yaml:"name" json:"name"`
		Namespace string `yaml:"namespace" json:"namespace"`
	} `yaml:"metadata" json:"metadata"`
	Spec interface{} `yaml:"spec" json:"spec"`
}

func PostConvertionHandlerFunc(c *gin.Context) {
	var param struct {
		Value string `json:"value"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		klog.Errorln("参数错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	if param.Value == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	timeNow := time.Now().Unix()
	saved := fmt.Sprintf("/tmp/%v.yaml", timeNow)
	saved2 := fmt.Sprintf("/tmp/%v-2.yaml", timeNow)
	err = ioutil.WriteFile(saved, []byte(param.Value),0777)
	if err != nil {
		klog.Errorln("保存文件错误", saved, err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "保存文件错误"))
		return
	}
	cmd := fmt.Sprintf(`cue import %s -l '"\(strings.ToCamel(kind))-\(strings.ToCamel(metadata.name))"' -o -`, saved)
	executor := exec.Command("bash","-c", cmd)
	var erro, outo bytes.Buffer
	executor.Stderr = &erro
	executor.Stdout = &outo
	err = executor.Run()
	if err != nil {
		klog.Errorln("执行命令错误", err.Error(), string(erro.Bytes()))
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, string(erro.Bytes())))
		return
	}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(param.Value)), 4096)
	yamls := ""
	for {
		out := generic{}
		err = decoder.Decode(&out)
		if err != nil {
			break
		}
		out.Metadata.Namespace = "context.namespace"
		out2, _ := yaml2.Marshal(out)
		yamls = fmt.Sprintf("%s\n---\n%s", yamls, string(out2))
	}
	err = ioutil.WriteFile(saved2, []byte(yamls),0777)
	if err != nil {
		klog.Errorln("保存文件错误", saved, err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "保存文件错误"))
		return
	}
	cmd2 := fmt.Sprintf(`cue import %s -l '"\(strings.ToCamel(kind))-\(strings.ToCamel(metadata.name))"' -o -`, saved2)
	executor2 := exec.Command("bash","-c", cmd2)
	var erro2, outo2 bytes.Buffer
	executor2.Stderr = &erro2
	executor2.Stdout = &outo2
	err = executor2.Run()
	if err != nil {
		klog.Errorln("执行命令错误", err.Error(), string(erro2.Bytes()))
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, string(erro2.Bytes())))
		return
	}
	c.JSON(200, utils.SuccessResponse(strings.ReplaceAll(outo2.String(), `"context.namespace"`, `context.namespace`)))
}