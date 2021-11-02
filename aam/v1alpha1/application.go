package v1alpha1

type Application struct {
	ApiVersion string
	Kind       string
	Metadata   struct {
		Name         string
		Version      string
		Description  string
		Keywords     []string
		Author       string
		Maintainers  []Maintainer
		Repositories []string
		Bugs         string
		Licenses     []License
		Annotations  map[string]string
	}
	Spec struct {
		Workloads    []Workload
		Exports      map[string][]string
		Dependencies []Dependency
		Userconfigs  map[string]interface{} `yaml:"userconfigs" json:"userconfigs"`
	}
}
