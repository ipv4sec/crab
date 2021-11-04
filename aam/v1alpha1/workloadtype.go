package v1alpha1

type WorkloadType struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct{
		Name string `yaml:"name"`
		Annotations map[string]string `yaml:"annotations"`
	} `yaml:"metadata"`
	Spec struct{
		Parameter string `yaml:"parameter"`
		Traits []string `yaml:"traits"`
	} `yaml:"spec"`
}
