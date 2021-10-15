package v1alpha1

type Manifest struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name        string      `yaml:"name" json:"name"`
		Annotations Annotations `yaml:"annotations" json:"annotations"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		Exports map[string][]string
		Dependencies []Dependency `yaml:"dependencies" json:"dependencies"`
		Configurations struct{} `yaml:"userconfigs" json:"userconfigs"`
		// TODO
		Components [] map[string]interface{} `json:"components"`
	} `yaml:"spec" json:"spec"`
}

