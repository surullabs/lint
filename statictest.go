package statictest

import (
	"go/build"
	"os"
	"reflect"
	"strings"
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

// Skipper is used to skip errors. Skip returns true if the error in
// str must be skipped.
type Skipper interface {
	Skip(str string) bool
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

// Skip returns a Checker that executes checkers and filters errors skipped by
// the provided skippers. If checker returns an Error instance the filters are
// applied on Errors.Errors. Each skipper is run in the order provided and a single
// skipper returning true will result in that error being skipped.
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

// Checker performs a static check of a single package. pkg must be a properly
// resolved import path. The behaviour for relative import paths is undefined.
// Use Apply(pkg, checker...) to perform the actual static checks.
type Checker interface {
	// Check performs a static check of all files in a package
	Check(pkg string) error
}

// CheckFunc is a function that implements Checker
type CheckFunc func(pkg string) error

// Check calls the CheckFunc with pkg
func (c CheckFunc) Check(pkg string) error { return c(pkg) }

// Apply applies the checkers to pkg. If pkg is a relative import path
// it will be resolved before being passed to the checker. checkers will be
// executed in the order provided. Errors are collected into an
// Error instance and returned.
func Apply(pkg string, checkers ...Checker) error {
	pkg, err := resolvePackage(pkg)
	if err != nil {
		return err
	}
	errs := &Error{}
	for _, checker := range checkers {
		name := reflect.TypeOf(checker).String()
		switch err := checker.Check(pkg).(type) {
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
}

func resolvePackage(dir string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	b, err := build.Import(dir, wd, build.FindOnly)
	if err != nil {
		return "", err
	}
	return b.ImportPath, nil
}
