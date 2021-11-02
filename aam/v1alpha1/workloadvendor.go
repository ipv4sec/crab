package v1alpha1

type WorkloadVendor struct {
	ApiVersion string
	Kind       string
	Metadata   struct{
		Name string
		Annotations map[string]string
	}
	Spec string
}
