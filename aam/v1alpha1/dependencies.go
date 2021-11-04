package v1alpha1

type Dependency struct {
	Instanceid string `json:"InstanceId"`
	Name string	`json:"name"`
	Version string	`json:"version"`
	Location string	`json:"location"`
	Items map[string][]string `json:"items"`
	EntryService string `json:"EntryService"`
}