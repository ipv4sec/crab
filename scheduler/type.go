package scheduler

type Component struct {
	ID   string
	Name string
	After string
	Deployment string

	HealthProbe map[string]string
	LoopNumber int
}
