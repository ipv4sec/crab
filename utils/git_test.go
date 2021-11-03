package utils

import "testing"

func TestInitRepo(t *testing.T) {
	err := InitRepo("/tmp/233", "https://github.com/ipv4sec/ez")
	if err != nil {
		if err.Error() != "already up-to-date" {
			t.Error(err.Error())
		}
	}
}
