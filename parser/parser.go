package parser

import "crab/aam/v1alpha1"

var (
	DIR_WORKLOAD_TYPE = "/Users/huanqiu/Desktop/uploads"
	DIR_WORKLOAD_VENDOR = "/Users/huanqiu/Desktop/uploads"
)

type VelaYaml struct {
	Name     string                 `json:"name"`
	Services map[string]interface{} `json:"services"`
}

type ManifestServiceOrigin struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name string `yaml:"name" json:"name"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		Components []Component `yaml:"components" json:"components"`
	} `yaml:"spec" json:"spec"`
}

type Component struct {
	Name       string `yaml:"name" json:"name"`
	Type       string `yaml:"type" json:"type"`
	Properties struct {
		Image         string          `yaml:"image" json:"image"`
		Port          int             `yaml:"port" json:"port"`
		Cmd           []string        `yaml:"cmd" json:"cmd"`
		Args          []string        `yaml:"args" json:"args"`
		Cpu           string          `yaml:"cpu" json:"cpu"`
		Init          string          `yaml:"init" json:"init"`
		Configs       []ConfigItem    `yaml:"configs" json:"configs"`
		After         string          `yaml:"after" json:"after"`
		Rootpwd       string          `yaml:"rootpwd" json:"rootpwd"`
		Storage       Storage         `yaml:"storage" json:"storage"`
		Authorization []Authorization `yaml:"authorization" json:"authorization"`
		Serviceentry  []ServiceEntry  `yaml:"serviceentry" json:"serviceentry"`
		Entry         Entry           `yaml:"entry" json:"entry"`
	} `yaml:"properties" json:"properties"`

	Traits []Trait `yaml:"traits" json:"traits"`
}
type Trait struct {
	Ttype      string            `yaml:"type" json:"ttype"`
	Properties map[string]string `yaml:"properties" json:"properties"`
}
type Export map[string][]string

type Annotations struct {
	Version     string   `yaml:"version" json:"version"`
	Description string   `yaml:"description" json:"description,omitempty"`
	Keywords    []string `yaml:"keywords" json:"keywords,omitempty"`
	Author      string   `yaml:"author" json:"author,omitempty"`
	Maintainers []struct {
		Name  string `yaml:"name" json:"name"`
		Email string `yaml:"email" json:"email"`
		Web   string `yaml:"web" json:"web"`
	} `yaml:"maintainers" json:"maintainers,omitempty"`
	Repositories []string `yaml:"repositories" json:"repositories,omitempty"`
	Bugs         string   `yaml:"bugs" json:"bugs,omitempty"`
	Licenses     []struct {
		Type string `yaml:"type" json:"type"`
	} `yaml:"licenses" json:"licenses"`
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
type Storage struct {
	Capacity string `yaml:"capacity" json:"capacity"`
	Path     string `yaml:"path" json:"path"`
}

//workload WebserviceVela
type WebserviceVela struct {
	Workload      string          `json:"workload"`
	Image         string          `json:"image"`
	Configs       []v1alpha1.ConfigItem    `json:"configs"`
	Storage       v1alpha1.Storage         `json:"storage"`
	Init          string          `json:"init,omitempty"`
	After         string          `json:"after,omitempty"`
	Port          int             `json:"port,omitempty"`
	Cmd           []string        `json:"cmd,omitempty"`
	Args          []string        `json:"args,omitempty"`
	Env           []v1alpha1.EnvItem       `json:"env,omitempty"`
	Traits        []struct {
		Type       string            `yaml:"type"`
		Properties map[string]interface{} `yaml:"properties"`
	}`json:"traits,omitempty"`
	Authorization []v1alpha1.Authorization `json:"authorization,omitempty"`
	Serviceentry  []v1alpha1.ServiceEntry  `json:"serviceentry,omitempty"`
	Namespace     string          `json:"namespace"`
	Type          string          `json:"type"`
	Entry         v1alpha1.Entry           `json:"entry,omitempty"`
}

//workload worker
type WorkerVela struct {
	Workload      string          `json:"workload"`
	Image         string          `json:"image"`
	Cmd           []string        `json:"cmd,omitempty"`
	Args          []string        `json:"args,omitempty"`
	Env           []v1alpha1.EnvItem       `json:"env,omitempty"`
	After         string          `json:"after,omitempty"`
	Init          string          `json:"init,omitempty"`
	Configs       []v1alpha1.ConfigItem    `json:"configs"`
	Storage       v1alpha1.Storage         `json:"storage"`
	Authorization []v1alpha1.Authorization `json:"authorization,omitempty"`
	Serviceentry  []v1alpha1.ServiceEntry  `json:"serviceentry,omitempty"`
	Namespace     string          `json:"namespace"`
	Type          string          `json:"type"`
	Entry         Entry           `json:"entry,omitempty"`
}

//workload redis
type RedisVela struct {
	Workload      string          `json:"workload"`
	After         string          `json:"after,omitempty"`
	Authorization []v1alpha1.Authorization `json:"authorization,omitempty"`
	Serviceentry  []v1alpha1.ServiceEntry  `json:"serviceentry,omitempty"`
	Namespace     string          `json:"namespace"`
	Type          string          `json:"type"`
	Entry         v1alpha1.Entry           `json:"entry,omitempty"`
}

//workload mysql
type MysqlVela struct {
	Workload      string          `json:"workload"`
	Rootpwd       string          `json:"rootpwd"`
	Storage       v1alpha1.Storage         `json:"storage"`
	Init          string          `json:"init,omitempty"`
	After         string          `json:"after,omitempty"`
	Authorization []v1alpha1.Authorization `json:"authorization,omitempty"`
	Serviceentry  []v1alpha1.ServiceEntry  `json:"serviceentry,omitempty"`
	Namespace     string          `json:"namespace"`
	Type          string          `json:"type"`
	Entry         v1alpha1.Entry           `json:"entry,omitempty"`
}

func (svc WebserviceVela) workload() string {
	return svc.Workload
}
func (svc WorkerVela) workload() string {
	return svc.Workload
}
func (svc MysqlVela) workload() string {
	return svc.Workload
}
func (svc RedisVela) workload() string {
	return svc.Workload
}

type Dependency struct {
	Instanceid   string              `json:"InstanceId"`
	Name         string              `json:"Name"`
	Location     string              `json:"Location"`
	Version      string              `json:"Version"`
	Uses         map[string][]string `json:"Uses"`
	EntryService string              `json:"EntryService"`
}

type DependencyVela struct {
	Instanceid   string              `json:"instanceid"`
	Name         string              `json:"name"`
	Location     string              `json:"location"`
	Version      string              `json:"version"`
	EntryService string              `json:"entryservice"`
	Resource     []DependencyUseItem `json:"resource"`
}

//解析后的依赖use
type DependencyUseItem struct {
	Uri     string   `json:"uri"`
	Actions []string `json:"actions"`
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

//返回的中间格式
type ParserData struct {
	Name string  `yaml:"name"`
	Init string `yaml:"init"`
	Workloads map[string]Workload `yaml:"workloads"`
}
type Workload struct {
	Parameter  string `yaml:"parameter"`
	Construct map[string]string `yaml:"construct"`
	Traits map[string]string `yaml:"traits"`
}

