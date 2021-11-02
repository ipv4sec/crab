package v1alpha1

type WorkloadType struct {
	ApiVersion string
	Kind       string
	Metadata   struct{
		Name string
		Annotations map[string]string
	}
	Spec struct{
		Parameter interface{}
		Traits []string
	}
}
