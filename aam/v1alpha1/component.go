package v1alpha1


type Component struct {
	Name string
	Properties map[string]interface{}
	Traits []Trait
}