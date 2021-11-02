package v1alpha1

type Dependency struct {
	Name     string
	Version  string
	Location string
	Items    map[string][]string
}
