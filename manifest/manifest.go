package manifest

import "crab/dependencies"


type VelaYaml struct {
	Name     string                 `json:"name"`
	Metadata Metadata               `json:"metadata"`
	Services map[string]interface{} `json:"services"`
}
type Metadata struct {
	Annotations Annotations `json:"annotations"`
}

type ManifestServiceOrigin struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name        string      `yaml:"name" json:"name"`
		//Annotations Annotations `yaml:"annotations" json:"annotations"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		Components []Component `yaml:"components" json:"components"`
	} `yaml:"spec" json:"spec"`
}

type Component struct {
	Name       string `yaml:"name" json:"name"`
	Type       string `yaml:"type" json:"type"`
	Properties struct {
		Image         string                       `yaml:"image" json:"image"`
		Port          int                          `yaml:"port" json:"port"`
		Cmd           []string                     `yaml:"cmd" json:"cmd"`
		Args          []string                     `yaml:"args" json:"args"`
		Cpu           string                       `yaml:"cpu" json:"cpu"`
		Init          string                       `yaml:"init" json:"init"`
		After         string                       `yaml:"after" json:"after"`
		Rootpwd       string                       `yaml:"rootpwd" json:"rootpwd"`
		Storage       Storage                      `yaml:"Storage" json:"storage"`
		Authorization []dependencies.Authorization `yaml:"authorization" json:"authorization"`
		Serviceentry  []dependencies.ServiceEntry  `yaml:"serviceentry" json:"serviceentry"`
		Entry         Entry                        `yaml:"entry" json:"entry"`
	} `yaml:"properties" json:"properties"`
	Traits []Trait `yaml:"traits" json:"traits"`
}
type Trait struct {
	Ttype string `yaml:"type" json:"ttype"`
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
	Capacity string `json:"capacity"`
	Path     string `json:"path"`
}

type Userconfig struct {
	Schema      string   `yaml:"$schema"`
	Id          string   `yaml:"$id"`
	Description string   `yaml:"description"`
	Type        string   `yaml:"type"`
	Properties  Property `yaml:"properties"`
	Required    []string `yaml:"required"`
}
type Param struct {
	Type       string `yaml:"type"`
	Properties map[string]struct {
		Items []struct {
			Type string `yaml:"type"`
		} `yaml:"item"`
		MinItems    int `yaml:"minItems"`
		UniqueItems int `yaml:"uniqueItems"`
	}
}
type Property struct {
	Param
	Required []string `yaml:"required"`
}

var WorkloadTypes = []string{"webservice", "worker", "redis", "mysql"}

//workload_vela
type WebserviceVela struct {
	Workload      string                     `json:"workload"`
	Image         string                     `json:"image"`
	Configs       []ConfigItem               `json:"configs"`
	Init          string                     `json:"init,omitempty"`
	After         string                     `json:"after,omitempty"`
	Port          int                        `json:"port,omitempty"`
	Cmd           []string                   `json:"cmd,omitempty"`
	Args          []string                   `json:"args,omitempty"`
	Env           []EnvItem                  `json:"env,omitempty"`
	Traits        []Trait                   `json:"traits"`
	Authorization []dependencies.Authorization `json:"authorization"`
	Serviceentry  []dependencies.ServiceEntry  `json:"serviceentry"`
	Namespace     string                     `json:"namespace"`
	Type          string                     `json:"type"`
	Entry         Entry                      `json:"entry"`
}

//workload worker
type WorkerVela struct {
	Workload      string                     `json:"workload"`
	Image         string                     `json:"image"`
	Cmd           []string                   `json:"cmd,omitempty"`
	Args          []string                   `json:"args,omitempty"`
	Env           []EnvItem                  `json:"env,omitempty"`
	After         string                     `json:"after,omitempty"`
	Init          string                     `json:"init,omitempty"`
	Configs       []ConfigItem               `json:"configs"`
	Storage       Storage                    `json:"storage"`
	Authorization []dependencies.Authorization `json:"authorization"`
	Serviceentry  []dependencies.ServiceEntry  `json:"serviceentry"`
	Namespace     string                     `json:"namespace"`
	Type          string                     `json:"type"`
	Entry         Entry                      `json:"entry"`
}

//workload redis
type RedisVela struct {
	Workload      string                     `json:"workload"`
	After         string                     `json:"after,omitempty"`
	Authorization []dependencies.Authorization `json:"authorization"`
	Serviceentry  []dependencies.ServiceEntry  `json:"serviceentry"`
	Namespace     string                     `json:"namespace"`
	Type          string                     `json:"type"`
	Entry         Entry                      `json:"entry"`
}

//workload mysql
type MysqlVela struct {
	Workload      string                     `json:"workload"`
	Rootpwd       string                     `json:"rootpwd"`
	Storage       Storage                    `json:"storage"`
	Init          string                     `json:"init,omitempty"`
	After         string                     `json:"after,omitempty"`
	Authorization []dependencies.Authorization `json:"authorization"`
	Serviceentry  []dependencies.ServiceEntry  `json:"serviceentry"`
	Namespace     string                     `json:"namespace"`
	Type          string                     `json:"type"`
	Entry         Entry                      `json:"entry"`
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
