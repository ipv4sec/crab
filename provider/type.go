package provider

type Dependency struct {
	Internal struct {
		Name         string
		ID           string `json:"InstanceId"`
		EntryService string
	}
	External struct {
		Name     string
		Location string
	}
}

func ConvertToDependency(param []struct {
	Name string `json:"name"`

	ID string `json:"id"`
	Location string `json:"location"`

	EntryService string
}) []Dependency {
	var val []Dependency
	for i := 0; i < len(param); i++ {
		if param[i].ID != "" {
			val = append(val, Dependency{
				Internal: struct {
					Name         string
					ID           string `json:"InstanceId"`
					EntryService string
				}{
					Name: param[i].Name,
					ID: param[i].ID,
					EntryService: param[i].EntryService,
				},
			})
		}
		if param[i].Location != "" {
			val = append(val, Dependency{
				External: struct {
					Name     string
					Location string
				}{Name: param[i].Name, Location: param[i].Location },
			})
		}
	}
	return val
}
