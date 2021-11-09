package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func Query(id, component string) error {
	if component == "" {
		return nil
	}
	res, err := HTTPClient.Get(fmt.Sprintf("http://island-api/status/%s/%s", id, component), nil)
	if err != nil {
		return fmt.Errorf("请求组件状态错误: %w", err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("读取组件状态返回错误: %w", err)
	}
	klog.Info("读取组件状态返回:", string(bodyBytes))
	var reply struct{
		Code int `json:"code"`
		Result interface{} `json:"result"`
	}
	err = json.Unmarshal(bodyBytes, &reply)
	if err != nil {
		return fmt.Errorf("组件状态返回序列化错误: %w", err)
	}
	if reply.Code != 0 {
		return fmt.Errorf("组件状态返回错误: %v %v", reply.Code, reply.Result)
	}
	if reply.Result.(float64) != 1 {
		return fmt.Errorf("组件状态未就绪: %v",reply.Result)
	}
	return nil
}