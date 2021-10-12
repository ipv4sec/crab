package provider

import (
	"bytes"
	"crab/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Yaml(manifest, uid, config string, dependencies []Dependency) (string, error) {

	requestByte, err := json.Marshal(struct {
		Manifest string `json:"content"`
		UID string `json:"instanceid"`
		Configuration string `json:"userconfig"`
		Dependencies []Dependency `json:"dependencies"`
	}{
		Manifest: manifest,
		UID: uid,
		Configuration: config,
		Dependencies: dependencies,
	})
	if err != nil {
		return "", fmt.Errorf("序列化错误: %w", err)
	}

	res, err := HTTPClient.Post("http://crab:5000/", bytes.NewReader(requestByte), nil)
	if err != nil {
		return "", fmt.Errorf("请求翻译器错误: %w", err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取翻译器返回错误: %w", err)
	}
	var reply utils.Reply
	err = json.Unmarshal(bodyBytes, &reply)
	if err != nil {
		return "", fmt.Errorf("TODO: %w", err)
	}
	if reply.Code != 0 {
		return "", fmt.Errorf("翻译器返回错误: %v", reply.Result)
	}
	return fmt.Sprintf("%v", reply.Result), nil
}
