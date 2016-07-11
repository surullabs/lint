package statictest

import "strings"

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
		n := make([]string, 0)
		for _, e := range err.Errors {
			if !skipper.Skip(e) {
				n = append(n, e)
			}
		}
		err.Errors = n
	}
	return err
}
