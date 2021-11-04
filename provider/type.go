package provider

type Dependency struct {
	ID           string `json:"InstanceId"`
	Uses         map[string][]string
	Name         string
	Version      string
	Location     string
	EntryService string
}

func ConvertToDependency(param []struct {
	ID string `json:"id"`
	Name string `json:"name"`
	EntryService string
}) []Dependency {
	var val []Dependency
	for i := 0; i < len(param); i++ {
		val = append(val, Dependency{
			ID:           param[i].ID,
			Uses:         nil,
			Name:         param[i].Name,
			Version:      "",
			Location:     "",
			EntryService: param[i].EntryService,
		})
	}
	return val
}
