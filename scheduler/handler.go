package scheduler

import (
	"context"
	"crab/cache"
	"crab/exec"
	"crab/parser"
	"crab/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

const (
	MessageQueueName = "crab:executor:deployments"
)

func PostDeploymentHandlerFunc(c *gin.Context)  {
	var param struct {
		Deployment string `json:"deployment"`
		ID string `json:"id"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}
	if param.ID == "" || param.Deployment == "" {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "参数错误"))
		return
	}

	var deployment parser.ParserData
	err = yaml.Unmarshal([]byte(param.Deployment), &deployment)
	if err != nil {
		klog.Errorln("反序列化失败:", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, err.Error()))
		return
	}
	for k, v := range deployment.Workloads {
		var parameters map[string] interface{}
		err = yaml.Unmarshal([]byte(v.Parameter), &parameters)
		if err != nil {
			klog.Errorln("反序列化工作负载参数失败:", err.Error())
			continue
		}
		component := Component{ID: param.ID, Name: k, LoopNumber: 0}
		after, ok := parameters["after"]
		if ok {
			component.After = after.(string)
		} else {
			component.Deployment = fmt.Sprintf("%v\n", deployment.Init)
		}
		for _, construct := range v.Construct {
			component.Deployment += fmt.Sprintf("---\n%s", construct)
		}
		for _, traits := range v.Traits {
			component.Deployment += fmt.Sprintf("---\n%s", traits)
		}
		component.HealthProbe = v.HealthProbe
		componentBytes, err := json.Marshal(component)
		if err != nil {
			klog.Errorln("序列化队列负载失败:", err.Error())
			continue
		}
		err = cache.Client.LPush(context.Background(), MessageQueueName, string(componentBytes)).Err()
		if err != nil {
			klog.Errorln("保存到队列失败", err.Error())
			continue
		}
	}
	c.JSON(200, utils.SuccessResponse("ok"))
}

func Consumption(){
	klog.Infoln("开始消费队列", time.Now().UTC())
	executor := exec.CommandExecutor{}
	for {
		time.Sleep(time.Second * 5)

		value, err := cache.Client.RPop(context.Background(), MessageQueueName).Result()
		if err != nil {
			if err != redis.Nil {
				panic(fmt.Errorf("消费队列出现错误: %w", err))
			}
			continue
		}
		var component Component
		err = json.Unmarshal([]byte(value), &component)
		if err != nil {
			klog.Infoln("反序列化错误", err.Error())
			continue
		}

		bytes, err := ioutil.ReadFile("assets/runtime/probe.yaml")
		if err != nil {
			panic(err)
		}
		probe := strings.ReplaceAll(strings.ReplaceAll(string(bytes), "$runtime",
			fmt.Sprintf("%s-%s", component.ID, component.Name)), "$command",
			component.HealthProbe["bash"])
		klog.Infoln("要执行的JOB内容为:", probe)
		saved := fmt.Sprintf("/tmp/%s_%s_job.yaml", component.ID, component.Name)
		err = ioutil.WriteFile(saved, []byte(probe),0777)
		if err != nil {
			klog.Errorln("保存文件错误", saved, err.Error())
			continue
		}
		command := fmt.Sprintf("/usr/local/bin/kubectl apply -f %s", saved)
		output, err := executor.ExecuteCommandWithCombinedOutput("bash", "-c", command)
		if err != nil {
			klog.Errorln("执行命令错误", err.Error())
		}

		err = HealthProbeStatus(fmt.Sprintf("profile-%s-%s", component.ID, component.Name))
		if err != nil {
			klog.Infoln("要执行的文件内容为:", component.Deployment)
			saved = fmt.Sprintf("/tmp/%s_%s_deployment.yaml", component.ID, component.Name)
			err = ioutil.WriteFile(saved, []byte(component.Deployment),0777)
			if err != nil {
				klog.Errorln("保存文件错误", saved, err.Error())
				continue
			}
			command = fmt.Sprintf("/usr/local/bin/kubectl apply -f %s", saved)
			output, err = executor.ExecuteCommandWithCombinedOutput("bash", "-c", command)
			if err != nil {
				klog.Errorln("执行命令错误", err.Error())
			}
			klog.Infoln("执行命令结果:", output)
			continue
		}
		// 未成功
		component.LoopNumber++
		v, err := json.Marshal(component)
		if err != nil {
			klog.Errorln("再次序列化错误", err.Error())
			continue
		}
		err = cache.Client.LPush(context.Background(), MessageQueueName, string(v)).Err()
		if err != nil {
			klog.Errorln("再次保存到队列失败", err.Error())
		}
	}
}