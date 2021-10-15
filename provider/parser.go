package provider

import (
	"bytes"
	"crab/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
	"mime/multipart"
)

func Yaml(manifest, uuid, domain string, config interface{}, dependencies []Dependency) (string, error) {
	v, err := json.Marshal(dependencies)
	if err != nil {
		return "", fmt.Errorf("序列化参数错误:%w", err)
	}

	requestByte := new(bytes.Buffer)
	w := multipart.NewWriter(requestByte)
	_ = w.WriteField("content", manifest)
	_ = w.WriteField("instanceid", uuid)
	_ = w.WriteField("userconfig", fmt.Sprintf("%v", config))
	_ = w.WriteField("dependencies", string(v))
	_ = w.WriteField("root-domain", domain)
	_ = w.Close()

	if err != nil {
		return "", fmt.Errorf("序列化错误: %w", err)
	}
	klog.Infoln("请求参数为:", requestByte.String())
	res, err := HTTPClient.Post("http://island-parser", requestByte, map[string][]string{
		"Content-Type": []string{w.FormDataContentType()},
	})
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
