package statictest

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os/exec"
	"reflect"
	"strings"
	"syscall"
)

// Error implements error and holds a list of failures.
type Error struct {
	Errors []string
}

func (e *Error) Error() string {
	return strings.Join(e.Errors, "\n")
}

// AsError returns nil if Errors is empty or the pointer to Error otherwise.
func (e *Error) AsError() error {
	if e.Errors == nil {
		return nil
	}
	return e
}

// Skipper is used to skip errors
type Skipper interface {
	Skip(string) bool
}

// StringSkipper implements Skipper and skips an error if Matcher(err, str) == true for
// any of Strings
type StringSkipper struct {
	Strings []string
	Matcher func(err, str string) bool
}

// Skip returns true if Matcher(check, str) == true for any of Strings.
func (s StringSkipper) Skip(check string) bool {
	for _, str := range s.Strings {
		if s.Matcher(check, str) {
			return true
		}
	}
	return false
}

func skip(check string, skippers []Skipper) bool {
	for _, s := range skippers {
		if s.Skip(check) {
			return true
		}
	}
	return false
}

// Skip executes checkers and filters errors skipped by the provided skippers. If checker returns
// an Error instance the filters are applied on Errors.Errors
func Skip(checker Checker, skippers ...Skipper) Checker {
	return CheckFunc(func(pkg string) error {
		switch err := checker.Check(pkg).(type) {
		case nil:
			return nil
		case *Error:
			var n []string
			for _, e := range err.Errors {
				if !skip(e, skippers) {
					n = append(n, e)
				}
			}
			err.Errors = n
			return err.AsError()
		default:
			if skip(err.Error(), skippers) {
				return nil
			}
			return err
		}
	})
}

// PackageDir returns the directory containing a package.
func PackageDir(path string) (string, error) {
	pkg, err := build.Import(path, ".", build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.Dir, nil
}

// InstallMissing runs go get importPath if bin cannot be found in the directories
// contained in the PATH environment variable.
func InstallMissing(bin, importPath string) error {
	if _, err := exec.LookPath(bin); err == nil {
		return nil
	}
	if data, err := exec.Command("go", "get", importPath).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install %s: %v: %s", importPath, err, string(data))
	}
	return nil
}

// ExecResult holds a status code, stdout and stderr for a single command execution.
type ExecResult struct {
	Code   int
	Stdout string
	Stderr string
}

// Exec executes cmd and results exit code, stdout and stderr of the result. It
// additionally returns and error if the status code is not 0.
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
	if err = cmd.Start(); err != nil {
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

// Checker performs a static check of a package
type Checker interface {
	// Check performs a static check of all files in a package
	Check(pkg string) error
}

// CheckFunc is a function that implements Checker
type CheckFunc func(pkg string) error

// Check calls the CheckFunc with pkg
func (c CheckFunc) Check(pkg string) error { return c(pkg) }

// Chain chains the provided checkers into a single checker. The returned checker
// will execute all checkers in the order provided. Errors are collected into an
// Error instance and returned.
func Chain(checkers ...Checker) Checker {
	return CheckFunc(func(dir string) error {
		errs := &Error{}
		for _, checker := range checkers {
			name := reflect.TypeOf(checker).String()
			switch err := checker.Check(dir).(type) {
			case nil:
				continue
			case *Error:
				for _, e := range err.Errors {
					errs.Errors = append(errs.Errors, name+": "+e)
				}
			default:
				errs.Errors = append(errs.Errors, name+": "+err.Error())
			}
		}
		return errs.AsError()
	})
}
