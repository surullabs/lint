package statictest

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

type Error struct {
	Errors []string
}

func (e *Error) Error() string {
	return strings.Join(e.Errors, "\n")
}

type Checker interface {
	Check(pkg string) error
}

type Skipper interface {
	Skip(string) bool
}

type stringSkipper []string

func (s stringSkipper) Skip(check string) bool {
	for _, str := range s {
		if str == check {
			return true
		}
	}
	return false
}

// SkipStrings returns a Skipper which will skip an error if the error is equal to
// any of strs.
func SkipStrings(strs ...string) Skipper { return stringSkipper(strs) }

// Skip removes errors skipped by skipper. err is returned unchanged if it is not
// of type *Error.
func Skip(err error, skipper Skipper) error {
	switch err := err.(type) {
	case *Error:
		var n []string
		for _, e := range err.Errors {
			if !skipper.Skip(e) {
				n = append(n, e)
			}
		}
		err.Errors = n
	}
	return err
}

func CheckCommand(name string, args ...string) error {
	data, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(data))
	}
	return nil
}

func PackageDir(path string) (string, error) {
	pkg, err := build.Import(path, ".", build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.Dir, nil
}

type ExecResult struct {
	Code   int
	Stdout string
	Stderr string
}

func Exec(cmd *exec.Cmd) (ExecResult, error) {
	res := ExecResult{Code: -1}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return res, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return res, err
	}
	if err := cmd.Start(); err != nil {
		return res, err
	}
	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		return res, fmt.Errorf("failed to read stdout: %v", err)
	}
	res.Stdout = string(data)
	if data, err = ioutil.ReadAll(stderr); err != nil {
		return res, fmt.Errorf("failed to read stderr: %v", err)
	}
	res.Stderr = string(data)
	err = cmd.Wait()
	if st, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		res.Code = st.ExitStatus()
	}
	return res, err
}
