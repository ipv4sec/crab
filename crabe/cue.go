package crabe

import (
	"github.com/blang/semver/v4"
	"github.com/pkg/errors"
	"strings"
)

func CheckCUEVersion() error {
	output, err := executor.ExecuteCommandWithCombinedOutput("cue", "version")
	if err != nil {
		return errors.Wrap(err, ErrCUEVersion)
	}
	vs := strings.Split(output, " ")

	if len(vs) != 4 {
		return errors.New(ErrCUEVersion)
	}

	version, err := semver.ParseTolerant(vs[2])
	if err != nil {
		return errors.Wrap(err, ErrCUEVersion)
	}

	if version.Minor < 4 {
		return errors.Errorf("%s (found %s)", ErrCUEVersion, vs[2])
	}
	return nil
}


func Export() {

}
