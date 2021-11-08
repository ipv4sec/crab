package provider

type Dependency struct {
	Name         string
	ID           string `json:"InstanceId"`
	EntryService string
	Location string
}
type Dependencies struct {
	Internal []Dependency
	External []Dependency
}

func ConvertToDependency(param []struct {
	Name string `json:"name"`

	ID string `json:"id"`
	Location string `json:"location"`

	EntryService string
}) Dependencies {
	var val Dependencies
	for i := 0; i < len(param); i++ {
		if param[i].ID != "" {
			val.Internal = append(val.Internal, Dependency{
					Name: param[i].Name,
					ID: param[i].ID,
					EntryService: param[i].EntryService,
				})
		}
		if param[i].Location != "" {
			val.External = append(val.External, Dependency{
				Name: param[i].Name,
				Location: param[i].Location,
			})
		}
	}
	return val
}
