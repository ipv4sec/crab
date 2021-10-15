package dependencies

import (
	"errors"
	"fmt"
)

type Dependency struct {
	Instanceid   string              `json:"instanceid"`
	Name         string              `json:"name"`
	Location     string              `json:"location"`
	Version      string              `json:"version"`
	Uses         map[string][]string `json:"uses"`
	EntryService string              `json:"entryservice"`
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

//解析use
func ApiParse(uses map[string][]string) ([]DependencyUseItem, error) {
	var err error
	rtn := make([]DependencyUseItem, 0)
	for k, v := range uses {
		count := 0
		actions := make([]string, 0)
		for _, option := range v {
			if option == "create" {
				actions = append(actions, "POST")
			}else if option == "read" {
				actions = append(actions, "GET", "HEAD", "OPTIONS")
			}else if option == "update" {
				actions = append(actions, "PUT", "PATCH")
			}else if option == "delete" {
				actions = append(actions, "DELETE")
			}else{
				return rtn, errors.New(fmt.Sprintf("依赖资源的操作类型(%s)不存在\n", option))
			}
			count++
		}
		if count == 0 {
			return rtn, errors.New("依赖资源的操作类型不能为空")
		}
		rtn = append(rtn, DependencyUseItem{k, actions})
	}
	return rtn,err
}
