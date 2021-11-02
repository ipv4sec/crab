package provider

//import (
//	"crab/aam/v1alpha1"
//	"crab/dependencies"
//	"encoding/json"
//	"fmt"
//	"gopkg.in/yaml.v3"
//	"io/ioutil"
//	"testing"
//)
//
//func TestYaml(t *testing.T) {
//	yamls := []string{
//		"examples/github/manifest.yaml",
//		"examples/harbor/manifest.yaml",
//		"examples/jenkins/manifest.yaml",
//		"examples/wstunnel/manifest.yaml",
//	}
//	for i := 0; i < len(yamls); i++ {
//		bytes, err := ioutil.ReadFile(yamls[i])
//		if err != nil {
//			panic(err)
//		}
//		var manifest v1alpha1.Manifest
//		err = yaml.Unmarshal(bytes, &manifest)
//		if err != nil {
//			panic(err)
//		}
//		value, err := json.Marshal(manifest)
//		if err != nil {
//			panic(err)
//		}
//		dep := `[{"name":"github","location":"https://www.github.com"}]`
//		var depe []dependencies.Dependency
//		err = yaml.Unmarshal([]byte(dep), &depe)
//		if err != nil {
//			panic(err)
//		}
//		// conf := `{"param1":"p1","param2":"p2","param3":{"param3_1":111,"param3_2":222},"param4":"p4"}"`
//
//		result, err := Yaml(string(value), "i233", "abc.com", struct {}{}, depe)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(result)
//	}
//}