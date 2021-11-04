package v1alpha1

type Properties map[string]interface{}


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