package lwcomponent

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

var (
	ErrNonExecutable error  = errors.New("component not executable")
	ErrRun           string = "unable to run component"
)

type Executer interface {
	Executable() bool

	Execute(args []string, envs ...string) (stdout string, stderr string, err error)

	ExecuteInline(args []string, envs ...string) (err error)
}

type executable struct {
	path string
}

func NewExecuable(name string, dir string) Executer {
	path := filepath.Join(dir, name)
	if runtime.GOOS == "windows" {
		path += ".exe"
	}

	return &executable{path: path}
}

func (e *executable) Executable() bool {
	return true
}

func (e *executable) Execute(args []string, envs ...string) (stdout string, stderr string, err error) {
	return execute(e.path, args, envs...)
}

func (e *executable) ExecuteInline(args []string, envs ...string) (err error) {
	return executeInline(e.path, args, envs...)
}

func execute(path string, args []string, envs ...string) (stdout string, stderr string, err error) {
	var outBuf, errBuf bytes.Buffer

	cmd := exec.Command(path, args...)

	cmd.Env = append(os.Environ(), envs...)

	cmd.Stdin = nil
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = run(cmd)

	stdout, stderr = outBuf.String(), errBuf.String()

	return
}

func executeInline(path string, args []string, envs ...string) error {
	cmd := exec.Command(path, args...)

	cmd.Env = append(os.Environ(), envs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return run(cmd)
}

func run(cmd *exec.Cmd) error {
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return &RunError{
				Err:      err,
				Message:  ErrRun,
				ExitCode: exitError.ExitCode(),
			}
		}
		return errors.Wrap(err, ErrRun)
	}

	return nil
}

type nonExecutable struct {
}

func (e *nonExecutable) Executable() bool {
	return false
}

func (e *nonExecutable) Execute(args []string, envs ...string) (stdout string, stderr string, err error) {
	return "", "", ErrNonExecutable
}

func (e *nonExecutable) ExecuteInline(args []string, envs ...string) (err error) {
	return ErrNonExecutable
}
