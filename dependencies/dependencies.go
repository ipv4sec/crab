package dependencies

type Dependency struct {
	Instanceid string              `json:"instanceid"`
	Name       string              `json:"name"`
	Location   string              `json:"location"`
	Version    string              `json:"version"`
	Uses       map[string][]string `json:"uses"`
}

type DependencyVela struct {
	Instanceid string              `json:"instanceid"`
	Name       string              `json:"name"`
	Location   string              `json:"location"`
	Version    string              `json:"version"`
	Resource   []DependencyUseItem `json:"resource"`
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

//解析use
func ApiParse(uses map[string][]string) []DependencyUseItem {
	rtn := make([]DependencyUseItem, 0)
	for k, v := range uses {
		rtn = append(rtn, DependencyUseItem{k, v})
	}
	return rtn
}
