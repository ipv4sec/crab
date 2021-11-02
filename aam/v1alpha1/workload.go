package v1alpha1

type Workload struct {
	Name       string
	Type       string
	Vendor     string
	Properties Properties
	Traits     []struct {
		Type       string
		Properties Properties
	}
}
