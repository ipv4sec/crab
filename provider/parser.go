package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func Yaml(manifest, uuid, entry string, config interface{}, dependencies Dependencies) (string, error, error) {
	v, err := json.Marshal(struct {
		Content string `json:"Content"`
		ID string `json:"InstanceId"`
		Configurations interface{} `json:"UserConfig"`
		Dependencies Dependencies `json:"Dependencies"`
		Host string `json:"Host"`
	}{
		Content: manifest,
		ID: uuid,
		Configurations: config,
		Dependencies: dependencies,
		Host: entry,
	})
	if err != nil {
		return "", nil, fmt.Errorf("序列化参数错误:%w", err)
	}

	klog.Infoln("请求参数为:", string(v))
	res, err := HTTPClient.Post("http://127.0.0.1:4000/", bytes.NewBuffer(v), nil)
	if err != nil {
		return "", nil, fmt.Errorf("请求翻译器错误: %w", err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, fmt.Errorf("读取翻译器返回错误: %w", err)
	}
	klog.Info("读取翻译器返回:", string(bodyBytes))
	var reply struct{
		Code int `json:"code"`
		Result interface{} `json:"result"`
	}
	err = json.Unmarshal(bodyBytes, &reply)
	if err != nil {
		return "", nil, fmt.Errorf("翻译器返回序列化错误: %w", err)
	}
	if reply.Code == 20101 {
		return "", fmt.Errorf("%v", reply.Result), nil
	}
	if reply.Code != 0 {
		return "", nil, fmt.Errorf("翻译器返回错误: %v", reply.Result)
	}
	return fmt.Sprintf("%v", reply.Result), nil, nil
}

func Template() (string, error) {
	res, err := HTTPClient.Get("http://127.0.0.1:4000/systemTemplate", nil)
	if err != nil {
		return "", fmt.Errorf("请求翻译器错误: %w", err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取翻译器返回错误: %w", err)
	}
	klog.Info("读取翻译器返回:", string(bodyBytes))
	var reply struct{
		Code int `json:"code"`
		Result interface{} `json:"result"`
	}
	err = json.Unmarshal(bodyBytes, &reply)
	if err != nil {
		return "", fmt.Errorf("翻译器返回序列化错误: %w", err)
	}
	if reply.Code != 0 {
		return "", fmt.Errorf("翻译器返回错误: %v", reply.Result)
	}
	return fmt.Sprintf("%v", reply.Result), nil
}
