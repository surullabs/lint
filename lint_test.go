package lint_test

import (
	"testing"

	"log"

	"fmt"
	"reflect"
	"runtime/debug"

	"strings"

	"github.com/surullabs/lint"
	"github.com/surullabs/lint/checkers"
	"github.com/surullabs/lint/dupl"
	"github.com/surullabs/lint/gofmt"
	"github.com/surullabs/lint/govet"
)

func TestLint(t *testing.T) {
	err := lint.Default.With(dupl.Check{Threshold: 25}).Check("./...")
	err = lint.Skip(err,
		// Ignore all errors from unused.go
		lint.RegexpMatch(`unused\.go`),
		// Ignore duplicates we're okay with.
		dupl.SkipTwo, dupl.Skip("golint.go:1,12"), dupl.Skip("errcheck.go:17,19"))

	if err != nil {
		t.Fatal(err)
	}
}

type checkFn func(pkgs ...string) error

func (c checkFn) Check(pkgs ...string) error { return c(pkgs...) }

func assert(t *testing.T, cond bool, msg string) {
	if cond {
		return
	}
	t.Fatal(msg, string(debug.Stack()))
}

var (
	expectRecursive = checkFn(func(args ...string) error {
		if !reflect.DeepEqual(args, []string{"./..."}) {
			return fmt.Errorf("expected [./...], got %v", args)
		}
		return nil
	})

	twoErrors = checkFn(func(...string) error {
		return checkers.Error("err1", "err2")
	})

	ungroupedError = checkFn(func(...string) error {
		return fmt.Errorf("ungrouped: %d", 1)
	})
)

func TestGroup(t *testing.T) {
	gcheck := func(fn ...lint.Checker) error {
		return lint.Group(fn).Check("./...")
	}

	// All checks pass
	err := gcheck(expectRecursive)
	assert(t, err == nil, fmt.Sprintf("%v", err))

	err = gcheck(twoErrors)
	assert(t,
		err != nil && err.Error() == "lint_test.checkFn: err1\nlint_test.checkFn: err2",
		fmt.Sprintf("%v", err))

	err = gcheck(ungroupedError)
	assert(t,
		err != nil && err.Error() == "lint_test.checkFn: ungrouped: 1",
		fmt.Sprintf("%v", err))

	// Grouped
	err = gcheck(twoErrors, ungroupedError)
	assert(t,
		err != nil && err.Error() == "lint_test.checkFn: err1\nlint_test.checkFn: err2\nlint_test.checkFn: ungrouped: 1",
		fmt.Sprintf("%v", err))
}

type skipFunc func(err string) bool

func (s skipFunc) Skip(err string) bool { return s(err) }

func errorIs(str string) skipFunc {
	return skipFunc(func(err string) bool { return str == err })
}

func scheck(c lint.Checker, skippers ...lint.Skipper) error {
	return lint.Skip(c.Check("./..."), skippers...)
}

func TestSkip(t *testing.T) {
	// All checks pass
	err := scheck(expectRecursive, errorIs("err1"))
	assert(t, err == nil, fmt.Sprintf("%v", err))

	err = scheck(twoErrors, errorIs("err1"))
	assert(t,
		err != nil && err.Error() == "err2",
		fmt.Sprintf("%v", err))

	// Skip ungrouped error
	err = scheck(ungroupedError, errorIs("ungrouped: 1"))
	assert(t, err == nil, fmt.Sprintf("%v", err))

	// Skip don't skip ungrouped error
	err = scheck(ungroupedError, errorIs("err1"))
	assert(t, err != nil && err.Error() == "ungrouped: 1", fmt.Sprintf("%v", err))

	// Grouped
	err = scheck(lint.Group{twoErrors, ungroupedError}, errorIs("lint_test.checkFn: err1"))
	assert(t,
		err != nil && err.Error() == "lint_test.checkFn: err2\nlint_test.checkFn: ungrouped: 1",
		fmt.Sprintf("%v", err))

}

func TestRegexpMatch(t *testing.T) {
	// regexp match
	skipRE := lint.RegexpMatch(`err1`, `ungrouped`)
	err := scheck(lint.Group{twoErrors, ungroupedError}, skipRE)
	assert(t,
		err != nil && err.Error() == "lint_test.checkFn: err2",
		fmt.Sprintf("%v", err))

	// panic on bad RE

	func() {
		defer func() {
			r := recover()
			if r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()
		_ = scheck(twoErrors, lint.RegexpMatch(`(unmatched paren`))
	}()
	assert(t,
		err != nil && strings.HasPrefix(err.Error(), "error parsing regexp"),
		fmt.Sprintf("%v", err))
}

func Example() {
	// Run the default set of linters
	err := lint.Default.Check("./...")

	// Ignore all errors from the file unused.go.
	//
	// This is intended as an example of how to skip errors and not a
	// recommendation that you skip these kinds of errors.
	err = lint.Skip(err, lint.RegexpMatch(
		`unused\.go:4:2: a blank import`,
		`unused\.go:7:7: don't use underscores in Go names`,
	))

	if err != nil {
		// Record lint failures.
		// Use t.Fatal(err) when running in a test
		log.Fatal(err)
	}
	// Output:
}

func Example_custom() {
	// Run a custom set of linters
	custom := lint.Group{
		gofmt.Check{}, govet.Check{},
	}
	if err := custom.Check("./..."); err != nil {
		// Record lint failures.
		// Use t.Fatal(err) when running in a test
		log.Fatal(err)
	}
	// Output:
}

func Example_packages() {
	custom := lint.Group{gofmt.Check{}, govet.Check{}}

	// Specify a list of packages instead of the wildcard
	if err := custom.Check(".", "./checkers"); err != nil {
		// Record lint failures.
		// Use t.Fatal(err) when running in a test
		log.Fatal(err)
	}
	// Output:
}
