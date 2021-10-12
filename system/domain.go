package system

import (
	"fmt"
	"io/ioutil"
)

func GetDomain() (string, error) {
	domain,err := ioutil.ReadFile("/etc/island-info/root-domain")
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return string(domain), nil
}
