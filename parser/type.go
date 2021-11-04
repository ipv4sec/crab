package parser

import "crab/aam/v1alpha1"

//workloadType的参数
type MysqlParam struct {
	Init string `json:"init"`
	Rootpwd string `json:"rootpwd"`
	Storage struct{
		Capacity string `json:"capacity"`
	} `json:"storage"`
	After string `json:"after"`
}
type WebserviceParam struct {
	Image string `json:"image"`
	Port int `json:"port"`
	Cmd []string `json:"cmd,omitempty"`
	Args []string `json:"args,omitempty"`
	Cpu string `json:"cpu,omitempty"`
	Env []v1alpha1.EnvItem `json:"env,omitempty"`
	After string `json:"after"`
}

type WorkerParam struct {
	Image string `json:"image"`
	Port int `json:"port"`
	Cmd []string `json:"cmd,omitempty"`
	Args []string `json:"args,omitempty"`
	Cpu string `json:"cpu,omitempty"`
	Env []v1alpha1.EnvItem `json:"env,omitempty"`
	After string `json:"after"`
}

type RedisParam struct {
	After string `json:"after"`
}

