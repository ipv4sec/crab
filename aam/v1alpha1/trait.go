package v1alpha1

type Trait struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   struct{
		Name string
		Annotations map[string]string
	} `yaml:"metadata"`
	Spec struct{
		Parameter interface{}
	}
}
