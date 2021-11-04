package v1alpha1

type Workload struct {
	Name       string	`yaml:"name"`
	Type       string	`yaml:"type"`
	Vendor     string	`yaml:"vendor"`
	Properties Properties `yaml:"properties"`
	Traits     []struct {
		Type       string	`yaml:"type"`
		Properties Properties `yaml:"properties"`
	} `yaml:"traits"`
}
