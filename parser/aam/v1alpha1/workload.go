package v1alpha1


type Workload struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
	Vendor string `json:"vendor" yaml:"vendor"`
	Properties map[string]interface{} `yaml:"properties" json:"properties"`
	Traits []Trait
}

type Trait struct {
	Type string `yaml:"type" json:"type"`
	Properties map[string]string `yaml:"properties" json:"properties"`
}