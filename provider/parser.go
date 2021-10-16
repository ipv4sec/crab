package provider

import (
	"bytes"
	dependency "crab/dependencies"
	"crab/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func Yaml(manifest, uuid, domain string, config interface{}, dependencies []dependency.Dependency) (string, error) {
	v, err := json.Marshal(struct {
		Manifest string `json:"content"`
		UUID string `json:"instanceid"`
		Configuration interface{} `json:"userconfig"`
		Dependencies []dependency.Dependency `json:"dependencies"`
		RootDomain string `json:"root-domain"`
	}{
		Manifest: manifest,
		UUID: uuid,
		Configuration: config,
		Dependencies: dependencies,
		RootDomain: domain,
	})
	if err != nil {
		return "", fmt.Errorf("序列化参数错误:%w", err)
	}

	klog.Infoln("请求参数为:", string(v))
	res, err := HTTPClient.Post("http://island-parser", bytes.NewBuffer(v), nil)
	if err != nil {
		return "", fmt.Errorf("请求翻译器错误: %w", err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取翻译器返回错误: %w", err)
	}
	klog.Info("读取翻译器返回:", string(bodyBytes))
	var reply utils.Reply
	err = json.Unmarshal(bodyBytes, &reply)
	if err != nil {
		return "", fmt.Errorf("翻译器返回序列化错误: %w", err)
	}
	if reply.Code != 0 {
		return "", fmt.Errorf("翻译器返回错误: %v", reply.Result)
	}
	return fmt.Sprintf("%v", reply.Result), nil
}
