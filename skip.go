package lint

import (
	"regexp"

	"github.com/surullabs/lint/checkers"
)

// Skipper is the interface that wraps the Skip method.
//
// Skip returns true if err is an error that must be ignored.
type Skipper interface {
	Skip(err string) bool
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

// Skip filters errors using skippers.  If err satisfies the below interface
//
//     type errors interface {
//     	Errors() []string
//     }
//
// the filters are applied to each string returned by Errors and then concatenated.
// If the error implements the above interface, it is guaranteed that any returned
// error will also implement the same.
//
// Skippers are run in the order provided and a single
// skipper returning true will result in that error being skipped.
func Skip(err error, skippers ...Skipper) error {
	switch serr := err.(type) {
	case nil:
		return nil
	case errors:
		var n []string
		errs := serr.Errors()
		for _, e := range errs {
			if !skip(e, skippers) {
				n = append(n, e)
			}
		}
		if len(n) == 0 {
			return nil
		}
		return checkers.Error(n...)
	default:
		if skip(serr.Error(), skippers) {
			return nil
		}
		return serr
	}
}

// RegexpMatch returns a Skipper that skips all errors which match
// any of the provided regular expression patterns. SkipRegexpMatch expects
// all patterns to be valid regexps and panics otherwise.
func RegexpMatch(regexps ...string) Skipper {
	return StringSkipper{
		Strings: regexps,
		Matcher: func(errstr, pattern string) bool {
			matched, err := regexp.MatchString(pattern, errstr)
			if err != nil {
				panic(err)
			}
			return matched
		},
	}
}
