package scheduler

import (
	"context"
	"crab/cache"
	"crab/exec"
	"crab/parser"
	"crab/provider"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"time"
)
type Result struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

type Param struct {
	Deploy string `json:"Deploy"`
	InstanceId string `json:"InstanceId"`
	Level string `json:"Level"`
	Operate string `json:"operate"`
}

func PostDeploymentHandlerFunc(c *gin.Context)  {
	var param Param
	err := c.ShouldBindJSON(&param)
	if err != nil {
		klog.Errorln(err)
		ret := Result{
			Code:   100011,
			Result: err.Error(),
		}
		c.JSON(200, ret)
		return
	}
	if param.InstanceId == "" ||  param.Deploy == "" {
		ret := Result{
			Code:   100010,
			Result: "参数不全",
		}
		c.JSON(200, ret)
		return
	}
	if param.Level == "" {
		param.Level = "develop"
	}
	if param.Operate == "" {
		param.Operate = "create"
	}
	//
	var deployment parser.ParserData
	err = yaml.Unmarshal([]byte(param.Deploy), &deployment)
	if err != nil {
		klog.Errorln("反序列化失败", err.Error())
	}
	for k, v := range deployment.Workloads {
		var parameters map[string] interface{}
		err = yaml.Unmarshal([]byte(v.Parameter), &parameters)
		if err != nil {
			klog.Errorln("反序列化Parameter失败:", err.Error())
			continue
		}
		component := Component{
			ID:    param.InstanceId,
			Name:  k,
		}
		after, ok := parameters["after"]
		if ok {
			component.After = after.(string)
		} else {
			component.Deployment = fmt.Sprintf("%v\n", deployment.Init)
		}
		for _, v2 := range v.Construct {
			component.Deployment += fmt.Sprintf("---\n%s", v2)
		}
		for _, v2 := range v.Traits {
			component.Deployment += fmt.Sprintf("---\n%s", v2)
		}
		componentBytes, err := json.Marshal(component)
		if err != nil {
			klog.Errorln("序列化component失败:", err.Error())
			continue
		}
		err = cache.Client.LPush(context.Background(), "crab:scheduler", string(componentBytes)).Err()
		if err != nil {
			klog.Errorln("保存到队列失败", err.Error())
		}
	}

	ret := Result{
		Code:   0,
		Result: "ok",
	}
	c.JSON(200, ret)
}

func Consumer(){
	klog.Infoln("开始消费队列", time.Now().UTC())
	executor := exec.CommandExecutor{}
	for {
		value, err := cache.Client.RPop(context.Background(), "crab:scheduler").Result()
		if err != nil {
			if err != redis.Nil {
				klog.Infoln("消费队列出现错误", err.Error())
				panic(err)
			}
			time.Sleep(time.Second * 5)
			continue
		}
		var component Component
		err = json.Unmarshal([]byte(value), &component)
		if err != nil {
			klog.Infoln("反序列化错误", err.Error())
			continue
		}

		err = provider.Query(component.ID, component.After)
		if err != nil {
			klog.Errorln("查询上一个组件状态错误", err.Error())
			err = cache.Client.LPush(context.Background(), "crab:scheduler", value).Err()
			if err != nil {
				klog.Errorln("保存到队列失败", err.Error())
			}
		}
		klog.Infoln("要执行的文件内容为:", component.Deployment)
		saved := fmt.Sprintf("/tmp/%s_%s.yaml", component.ID, component.Name)
		err = ioutil.WriteFile(saved, []byte(component.Deployment),0777)
		if err != nil {
			klog.Errorln("保存文件错误", saved, err.Error())
			continue
		}
		command := fmt.Sprintf("/usr/local/bin/kubectl apply -f %s", saved)
		output, err := executor.ExecuteCommandWithCombinedOutput("bash", "-c", command)
		if err != nil {
			klog.Errorln("执行命令错误", err.Error())
		}
		klog.Infoln("执行命令结果:", output)
	}
}