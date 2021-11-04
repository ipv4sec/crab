package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func Exec(id, value string) error {
	v, err := json.Marshal(struct {
		ID string `json:"InstanceId"`
		Content string `json:"Yaml"`
	}{
		Content: value,
		ID: id,
	})
	if err != nil {
		return fmt.Errorf("序列化参数错误:%w", err)
	}

	klog.Infoln("请求参数为:", string(v))
	res, err := HTTPClient.Post("http://island-scheduler", bytes.NewBuffer(v), nil)
	if err != nil {
		return fmt.Errorf("请求调度器错误: %w", err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("读取调度器返回错误: %w", err)
	}
	klog.Info("读取调度器返回:", string(bodyBytes))
	var reply struct{
		Code int `json:"code"`
		Result string `json:"result"`
	}
	err = json.Unmarshal(bodyBytes, &reply)
	if err != nil {
		return fmt.Errorf("调度器返回序列化错误: %w", err)
	}
	if reply.Code != 0 || reply.Result != "ok"{
		return fmt.Errorf("调度器返回错误: %v %v", reply.Code, reply.Result)
	}
	return nil
}