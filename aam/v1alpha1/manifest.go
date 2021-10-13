package v1alpha1

import "encoding/json"

type Manifest struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name        string      `yaml:"name" json:"name"`
		Annotations Annotations `yaml:"annotations" json:"annotations"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		Exports map[string][]string
		Dependencies []Dependency
		Configurations interface{} `yaml:"userconfigs" json:"userconfigs"`
		Components []json.RawMessage `json:"components"`
	} `yaml:"spec" json:"spec"`
}

