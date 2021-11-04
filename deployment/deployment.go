package deployment

type Deployment struct {
	ID string `json:"id"`
	Domain string `json:"domain"`
	Configurations interface{} `json:"userconfigs"`
	Dependencies string `json:"dependencies"`
}