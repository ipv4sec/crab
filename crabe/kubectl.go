package crabe

import (
	"encoding/json"
	"github.com/blang/semver/v4"
	"github.com/pkg/errors"
)

func CheckKubectlVersion() error {
	output, err := executor.ExecuteCommandWithCombinedOutput("kubectl", "version", "--client", "-o", "json")
	if err != nil {
		return errors.Wrap(err, ErrKubectlVersion)
	}

	var info struct {
		ClientVersion struct {
			GitVersion string
		}
	}

	if err = json.Unmarshal([]byte(output), &info); err != nil {
		return errors.Wrap(err, ErrKubectlVersion)
	}

	version, err := semver.ParseTolerant(info.ClientVersion.GitVersion)
	if err != nil {
		return errors.Wrap(err, ErrKubectlVersion)
	}

	if version.Major != 1 || version.Minor < 10 {
		return errors.Errorf("%s (found %s)", ErrKubectlVersion, info.ClientVersion.GitVersion)
	}
	return nil
}

