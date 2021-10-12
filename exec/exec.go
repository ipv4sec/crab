/*
Copyright 2016 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Executor is the main interface for all the exec commands
type Executor interface {
	ExecuteCommand(command string, arg ...string) error
	ExecuteCommandWithEnv(env []string, command string, arg ...string) error
	ExecuteCommandWithOutput(command string, arg ...string) (string, error)
	ExecuteCommandWithCombinedOutput(command string, arg ...string) (string, error)
	ExecuteCommandWithOutputFile(command, outfileArg string, arg ...string) (string, error)
	ExecuteCommandWithOutputFileTimeout(timeout time.Duration, command, outfileArg string, arg ...string) (string, error)
	ExecuteCommandWithTimeout(timeout time.Duration, command string, arg ...string) (string, error)
}

// CommandExecutor is the type of the Executor
type CommandExecutor struct {
}

// ExecuteCommand starts a process and wait for its completion
func (c *CommandExecutor) ExecuteCommand(command string, arg ...string) error {
	return c.ExecuteCommandWithEnv([]string{}, command, arg...)
}

// ExecuteCommandWithEnv starts a process with env variables and wait for its completion
func (*CommandExecutor) ExecuteCommandWithEnv(env []string, command string, arg ...string) error {
	cmd, _, _, err := startCommand(env, command, arg...)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// ExecuteCommandWithTimeout starts a process and wait for its completion with timeout.
func (*CommandExecutor) ExecuteCommandWithTimeout(timeout time.Duration, command string, arg ...string) (string, error) {
	// #nosec G204 Rook controls the input to the exec arguments
	cmd := exec.Command(command, arg...)

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	if err := cmd.Start(); err != nil {
		return "", err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	interruptSent := false
	for {
		select {
		case <-time.After(timeout):
			if interruptSent {
				var e error
				if err := cmd.Process.Kill(); err != nil {
					e = fmt.Errorf("timeout waiting for the command %s to return after interrupt signal was sent. Tried to kill the process but that failed: %v", command, err)
				} else {
					e = fmt.Errorf("timeout waiting for the command %s to return", command)
				}
				return strings.TrimSpace(b.String()), e
			}

			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				// kill signal will be sent next loop
			}
			interruptSent = true
		case err := <-done:
			if err != nil {
				return strings.TrimSpace(b.String()), err
			}
			if interruptSent {
				return strings.TrimSpace(b.String()), fmt.Errorf("timeout waiting for the command %s to return", command)
			}
			return strings.TrimSpace(b.String()), nil
		}
	}
}

// ExecuteCommandWithOutput executes a command with output
func (*CommandExecutor) ExecuteCommandWithOutput(command string, arg ...string) (string, error) {
	// #nosec G204 Rook controls the input to the exec arguments
	cmd := exec.Command(command, arg...)
	return runCommandWithOutput(cmd, false)
}

// ExecuteCommandWithCombinedOutput executes a command with combined output
func (*CommandExecutor) ExecuteCommandWithCombinedOutput(command string, arg ...string) (string, error) {
	// #nosec G204 Rook controls the input to the exec arguments
	cmd := exec.Command(command, arg...)
	return runCommandWithOutput(cmd, true)
}

// ExecuteCommandWithOutputFileTimeout Same as ExecuteCommandWithOutputFile but with a timeout limit.
// #nosec G307 Calling defer to close the file without checking the error return is not a risk for a simple file open and close
func (*CommandExecutor) ExecuteCommandWithOutputFileTimeout(timeout time.Duration,
	command, outfileArg string, arg ...string) (string, error) {

	outFile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", errors.Wrap(err, "failed to open output file")
	}
	defer outFile.Close()
	defer os.Remove(outFile.Name())

	arg = append(arg, outfileArg, outFile.Name())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// #nosec G204 Rook controls the input to the exec arguments
	cmd := exec.CommandContext(ctx, command, arg...)
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		cmdOut = []byte(fmt.Sprintf("%s. %s", string(cmdOut), assertErrorType(err)))
	}

	// if there was anything that went to stdout/stderr then log it, even before
	// we return an error
	if string(cmdOut) != "" {
	}

	if ctx.Err() == context.DeadlineExceeded {
		return string(cmdOut), ctx.Err()
	}

	if err != nil {
		return string(cmdOut), &CephCLIError{err: err, output: string(cmdOut)}
	}

	fileOut, err := ioutil.ReadAll(outFile)
	if err := outFile.Close(); err != nil {
		return "", err
	}
	return string(fileOut), err
}

// ExecuteCommandWithOutputFile executes a command with output on a file
// #nosec G307 Calling defer to close the file without checking the error return is not a risk for a simple file open and close
func (*CommandExecutor) ExecuteCommandWithOutputFile(command, outfileArg string, arg ...string) (string, error) {

	// create a temporary file to serve as the output file for the command to be run and ensure
	// it is cleaned up after this function is done
	outFile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", errors.Wrap(err, "failed to open output file")
	}
	defer outFile.Close()
	defer os.Remove(outFile.Name())

	// append the output file argument to the list or args
	arg = append(arg, outfileArg, outFile.Name())
	// #nosec G204 Rook controls the input to the exec arguments
	cmd := exec.Command(command, arg...)
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		cmdOut = []byte(fmt.Sprintf("%s. %s", string(cmdOut), assertErrorType(err)))
	}
	// if there was anything that went to stdout/stderr then log it, even before we return an error
	if string(cmdOut) != "" {
	}
	if err != nil {
		return string(cmdOut), &CephCLIError{err: err, output: string(cmdOut)}
	}

	// read the entire output file and return that to the caller
	fileOut, err := ioutil.ReadAll(outFile)
	if err := outFile.Close(); err != nil {
		return "", err
	}
	return string(fileOut), err
}

func startCommand(env []string, command string, arg ...string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {

	// #nosec G204 Rook controls the input to the exec arguments
	cmd := exec.Command(command, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
	}

	if len(env) > 0 {
		cmd.Env = env
	}

	err = cmd.Start()

	return cmd, stdout, stderr, err
}

func runCommandWithOutput(cmd *exec.Cmd, combinedOutput bool) (string, error) {
	var output []byte
	var err error
	var out string

	if combinedOutput {
		output, err = cmd.CombinedOutput()
	} else {
		output, err = cmd.Output()
		if err != nil {
			output = []byte(fmt.Sprintf("%s. %s", string(output), assertErrorType(err)))
		}
	}

	out = strings.TrimSpace(string(output))

	if err != nil {
		return out, err
	}

	return out, nil
}

func assertErrorType(err error) string {
	switch errType := err.(type) {
	case *exec.ExitError:
		return string(errType.Stderr)
	case *exec.Error:
		return errType.Error()
	}

	return ""
}

// ExtractExitCode attempts to get the exit code from the error returned by an Executor function.
// This should also work for any errors returned by the golang os/exec package.
func ExtractExitCode(err error) (int, error) {
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode(), nil
	}
	return 0, errors.Errorf("error %#v is not an ExitError", err)
}
