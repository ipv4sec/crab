package v1alpha1


type Metadata struct {
	Name string `yaml:"name" json:"name"`
	Version	string `yaml:"version" json:"version"`
	Description string `yaml:"description" json:"description"`
	Keywords []string `yaml:"keywords" json:"keywords"`
	Author string `yaml:"author" json:"author"`
	Maintainers []Maintainer `yaml:"maintainers" json:"maintainers"`
	Repositories []string `yaml:"repositories" json:"repositories"`
	Bugs string `yaml:"bugs" json:"bugs"`
	Licenses []License `yaml:"license" json:"license"`
	Annotations map[string]string `yaml:"annotations" json:"annotations"`
}

type Maintainer struct {
	Email string `yaml:"email" json:"email"`
	Name string `yaml:"name" json:"name"`
	Web string `yaml:"web" json:"web"`
}

type License struct {
	Type string `yaml:"type" json:"type"`
	Url string `yaml:"url" json:"url"`
}