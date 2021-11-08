package parser

import "crab/aam/v1alpha1"

type ContextObj struct {
	AppName       string `json:"appName"`
	ComponentName string `json:"componentName"`
	Namespace     string `json:"namespace"`
}
type Result struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}
type Dependency struct {
	Instanceid   string              `json:"InstanceId"`
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Location     string              `json:"location"`
	Items        map[string][]string `json:"items"`
	EntryService string              `json:"EntryService"`
}

//验证type,vendor返回的数据
type WorkloadParam struct {
	Parameter map[string]interface{} `json:"parameter"`
	Type      string                 `json:"type"`
	Vendor    string                 `json:"vendor"`
	VendorCue string                 `json:"vendorCue"`
	Traits    []string               `json:"traits"`
}
type VelaYaml struct {
	Name     string                 `json:"name"`
	Services map[string]interface{} `json:"services"`
}

//返回的中间格式
type ParserData struct {
	Name      string              `yaml:"name"`
	Init      string              `yaml:"init"`
	Workloads map[string]Workload `yaml:"workloads"`
}
type Workload struct {
	Parameter string            `yaml:"parameter"`
	Construct map[string]string `yaml:"construct"`
	Traits    map[string]string `yaml:"traits"`
}
type ConfigItem struct {
	Path    string               `yaml:"path" json:"path"`
	SubPath string               `yaml:"subPath" json:"subPath,omitempty"`
	Data    []ConfigItemDataItem `yaml:"data" json:"data"`
}
type ConfigItemDataItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Storage struct {
	Capacity string `yaml:"capacity" json:"capacity"`
	Path     string `yaml:"path" json:"path"`
}

//内部应用授权
type Authorization struct {
	Namespace string              `json:"namespace"`
	Service   string              `json:"service"`
	Resources []DependencyUseItem `json:"resources,omitempty"`
}

//外部应用授权
type ServiceEntry struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

//解析后的依赖use
type DependencyUseItem struct {
	Uri     string   `json:"uri"`
	Actions []string `json:"actions"`
}
type EnvItem struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	valueFrom map[string]struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	}
}
type Entry struct {
	Host string   `json:"host"`
	Path []string `json:"path"`
}
type DependencyVela struct {
	Instanceid   string              `json:"instanceid"`
	Name         string              `json:"name"`
	Location     string              `json:"location"`
	EntryService string              `json:"entryservice"`
	Resource     []DependencyUseItem `json:"resource"`
}
type Trait struct {
	Type       string              `yaml:"type"`
	Properties v1alpha1.Properties `yaml:"properties"`
}

const (
	ErrBadRequest     = 10201
	ErrInternalServer = 10202
)
