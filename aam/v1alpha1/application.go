package v1alpha1

type Application struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name         string `yaml:"name" json:"name"`
		Version      string `yaml:"version" json:"version"`
		Description  string `yaml:"description" json:"description"`
		Keywords     []string `yaml:"keywords" json:"keywords"`
		Author       string `yaml:"author" json:"author"`
		Maintainers  []Maintainer `yaml:"maintainers" json:"maintainers"`
		Repositories []string `yaml:"repositories" json:"repositories"`
		Bugs         string `yaml:"bugs" json:"bugs"`
		Licenses     []License `yaml:"licenses" json:"licenses"`
		Annotations  map[string]string `yaml:"annotations" json:"annotations"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		Workloads    []Workload `yaml:"workloads" json:"workloads"`
		Exports      map[string][]string
		Dependencies []Dependency
		Userconfigs  map[string]interface{} `yaml:"userconfigs" json:"userconfigs"`
	} `yaml:"spec" json:"spec"`
}
