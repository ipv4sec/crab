package v1alpha1

type Annotations struct {
	Version     string   `yaml:"version" json:"version"`
	Description string   `yaml:"description" json:"description"`
	Keywords    []string `yaml:"keywords" json:"keywords"`
	Author      string   `yaml:"author" json:"author"`
	Maintainers []struct {
		Name  string `yaml:"name" json:"name"`
		Email string `yaml:"email" json:"email"`
		Web   string `yaml:"web" json:"web"`
	} `yaml:"maintainers" json:"maintainers"`
	Repositories []string `yaml:"repositories" json:"repositories"`
	Bugs         string   `yaml:"bugs" json:"bugs"`
	Licenses     []struct {
		Type string `yaml:"type" json:"type"`
	} `yaml:"licenses" json:"licenses"`
}
