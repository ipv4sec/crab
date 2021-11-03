package v1alpha1


type Manifest struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind string `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec struct {
		Workloads []Workload `yaml:"workloads" json:"workloads"`
	} `yaml:"spec" json:"spec"`
}