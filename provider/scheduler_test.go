package provider

import (
	"crab/aam/v1alpha1"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestExec(t *testing.T) {
	yamls := []string{
		"examples/github/manifest.yaml",
		"examples/harbor/manifest.yaml",
		"examples/jenkins/manifest.yaml",
		"examples/wstunnel/manifest.yaml",
	}
	for i := 0; i < len(yamls); i++ {
		bytes, err := ioutil.ReadFile(yamls[i])
		if err != nil {
			panic(err)
		}
		var manifest v1alpha1.Application
		err = yaml.Unmarshal(bytes, &manifest)
		if err != nil {
			panic(err)
		}
		value, err := json.Marshal(manifest)
		if err != nil {
			panic(err)
		}
		var dep []Dependency
		err = yaml.Unmarshal([]byte(`[{"name":"github","location":"https://www.github.com"}]`), &dep)
		if err != nil {
			panic(err)
		}
		conf := `{"param1":"p1","param2":"p2","param3":{"param3_1":111,"param3_2":222},"param4":"p4"}"`

		result, err := Yaml(string(value), "i233", "abc.com", conf, dep)
		if err != nil {
			panic(err)
		}
		err = Exec("i233", result)
		if err != nil {
			panic(err)
		}
	}
}