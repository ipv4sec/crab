package scheduler

import (
	"crab/cache"
	"crab/parser"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"context"
	"time"
	"os/exec"
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

func PostDeployHandlerFunc(c *gin.Context)  {
	var param Param
	var err error
	err = c.ShouldBind(&param)
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
	str,err := json.Marshal(param)
	if err != nil {
		ret := Result{
			Code:   100012,
			Result: "序列化失败",
		}
		c.JSON(200, ret)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	err = cache.Client.LPush(ctx, "crab", string(str)).Err()
	if err != nil {
		// TODO
		ret := Result{
			Code:   100010,
			Result: err.Error(),
		}
		klog.Errorln("err")
		c.JSON(200, ret)
		return
	}
	ret := Result{
		Code:   0,
		Result: "开始部署",
	}
	c.JSON(200, ret)
}

func Consumer(){
	klog.Infoln("开始消费队列", time.Now().UTC())
	for {
		value, err := cache.Client.LPop(context.Background(), "crab").Result()
		if err != nil {
			if err != redis.Nil {
				klog.Infoln("消费队列出现错误", err.Error())
				panic(err)
			}
			time.Sleep(time.Second * 5)
			continue
		}
		var m Param
		err = json.Unmarshal([]byte(value), &m)
		if err != nil {
			klog.Infoln("反序列化错误", err.Error())
			continue
		}
		var deploy parser.ParserData
		err = yaml.Unmarshal([]byte(m.Deploy), &deploy)
		if err != nil {
			klog.Errorln("deploy反序列化失败")
			continue
		}
		fmt.Println(deploy.Init)
		if m.Operate == "create" || m.Operate == "update" {
			saved := fmt.Sprintf("tmp/%s_%s.yaml", m.Operate, m.InstanceId)
			err = ioutil.WriteFile(saved, []byte(deploy.Init),0777)
			if err != nil {
				klog.Errorln("保存文件错误", saved, err.Error())
				continue
			}
			command := fmt.Sprintf("/usr/local/bin/kubectl apply -f %s", saved)
			cmd := exec.Command("bash", "-c", command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				klog.Errorln("执行命令错误", err.Error())
			}
			fmt.Println(string(output))
		}
		time.Sleep(time.Second * 2)
	}
}