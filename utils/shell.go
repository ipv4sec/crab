package utils

import (
	"os"
	"os/exec"
	"strings"
)

func Exec(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "KUBECONFIG=/etc/kubernetes/admin.conf")
	cmd.Env = append(cmd.Env, "HOME=/var/island")
	cmd.Dir = "/var/island"
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}
