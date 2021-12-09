package crabe

import (
	"crab/exec"
	"encoding/json"
	"fmt"
	"os"
)

var executor = exec.CommandExecutor{}
const ErrKubectlVersion = "kubectl version 1.10 or greater is required"
const ErrCUEVersion = "kubectl version 1.10 or greater is required"

func Prepare(workspace string) {
	err := os.Mkdir(fmt.Sprintf("/tmp/crab/%s", workspace), os.ModePerm)
	if err!=nil{
		fmt.Println(err)
	}
}

func Linking(template string, runtime interface{}) string {
	var err error
	err = CheckCUEVersion()
	if err != nil {}

	var t Template
	err = json.Unmarshal([]byte(template), &t)
	if err != nil {}
	// 保存到工作区

	output, err := executor.ExecuteCommandWithCombinedOutput("cue", "export", "-f", "path")
	if err != nil {
		return ""
	}
	var result map[string]map[string]interface{}
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {}

	// 合并

	return "/tmp/crab/%s/deployment.yaml"
}

func Deployment(yamlPath, serverUrl string) {
	var err error
	err = CheckKubectlVersion()
	if err != nil {}
	args := []string{"apply", "-f", yamlPath}
	if serverUrl != "" {
		args = append(args, "-s", serverUrl)
	}
	output, err := executor.ExecuteCommandWithCombinedOutput("kubectl", args...)
	if err != nil {
		return
	}
	fmt.Println(output)
}