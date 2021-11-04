package v1alpha1

//type Properties map[string]interface{}
type Properties struct {
	Image         string          `yaml:"image" json:"image,omitempty"`
	Port          int             `yaml:"port" json:"port,omitempty"`
	Cmd           []string        `yaml:"cmd" json:"cmd,omitempty"`
	Args          []string        `yaml:"args" json:"args,omitempty"`
	Cpu           string          `yaml:"cpu" json:"cpu,omitempty"`
	Init          string          `yaml:"init" json:"init,omitempty"`
	Configs       []ConfigItem    `yaml:"configs" json:"configs,omitempty"`
	Env       	  []EnvItem       `yaml:"env" json:"env,omitempty"`
	After         string          `yaml:"after" json:"after"`
	Rootpwd       string          `yaml:"rootpwd" json:"rootpwd"`
	Storage       Storage         `yaml:"storage" json:"storage"`
	Authorization []Authorization `yaml:"authorization" json:"authorization,omitempty"`
	Serviceentry  []ServiceEntry  `yaml:"serviceentry" json:"serviceentry,omitempty"`
	Entry         Entry           `yaml:"entry" json:"entry,omitempty"`
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
	Version      string              `json:"version"`
	EntryService string              `json:"entryservice"`
	Resource     []DependencyUseItem `json:"resource"`
}