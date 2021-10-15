package v1alpha1

type Dependency struct {
	Name     string `yaml:"name" json:"name"`
	Version  string `yaml:"version" json:"version"`
	Location string `yaml:"location" json:"location"`
	Uses     map[string][]string `yaml:"uses" json:"uses"`
}
