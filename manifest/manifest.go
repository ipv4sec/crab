package manifest

import "crab/dependencies"

type Manifest struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name        string      `yaml:"name" json:"name"`
		Annotations Annotations `yaml:"annotations" json:"annotations"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		Components []Component `yaml:"components" json:"components"`
	} `yaml:"spec" json:"spec"`
}

type VelaYaml struct {
	Name     string                 `yaml:"name"`
	Services map[string]interface{} `yaml:"services"`
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
		Traits        []string                     `yaml:"traits" json:"traits"`
		After         string                       `yaml:"after" json:"after"`
		Rootpwd       string                       `yaml:"rootpwd" json:"rootpwd"`
		Storage       Storage                      `yaml:"Storage" json:"storage"`
		Authorization []dependencies.Authorization `yaml:"authorization" json:"authorization"`
		Serviceentry  []dependencies.ServiceEntry  `yaml:"serviceentry" json:"serviceentry"`
		Entry         Entry                        `yaml:"entry" json:"entry"`
	} `yaml:"properties" json:"properties"`
}
type Export map[string][]string

type Annotations struct {
	Version     string   `yaml:"version" json:"version"`
	Description string   `yaml:"description" json:"description"`
	Keywords    []string `yaml:"keywords" json:"keywords"`
	Author      string   `yaml:"author" json:"author"`
	Maintainers []struct {
		Name  string `yaml:"name" json:"name"`
		Email string `yaml:"email" json:"email"`
		Web   string `yaml:"web" json:"web"`
	} `yaml:"maintainers" json:"maintainers"`
	Repositories []string `yaml:"repositories" json:"repositories"`
	Bugs         string   `yaml:"bugs" json:"bugs"`
	Licenses     []struct {
		Type string `yaml:"type" json:"type"`
	} `yaml:"licenses" json:"licenses"`
}
type ConfigItem struct {
	Path    string               `yaml:"path" json:"path"`
	SubPath string               `yaml:"subPath" json:"subPath"`
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
