package deployment

type Deployment struct {
	ID string `json:"id"`
	Domain string `json:"domain"`
	Configurations interface{} `json:"userconfigs" yaml:"userconfigs"`
	Dependencies []struct{
		Name string `json:"name"`

		ID string `json:"id"`
		Location string `json:"location"`

		EntryService string
	} `json:"dependencies"`
}