package provider

type Dependency struct {
	Name string `json:"name"`

	UID  string `json:"instanceid"`
	Location string `json:"location"`
}