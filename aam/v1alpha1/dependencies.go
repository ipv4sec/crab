package v1alpha1

type Dependency struct {
	Name string	`json:"name"`
	Version string	`json:"version"`
	Location string	`json:"location"`
	Items map[string][]string `json:"items"`
}