package lint_test

import (
	"testing"

	"log"

	"github.com/surullabs/lint"
	"github.com/surullabs/lint/dupl"
	"github.com/surullabs/lint/gofmt"
	"github.com/surullabs/lint/golint"
	"github.com/surullabs/lint/gosimple"
	"github.com/surullabs/lint/gostaticcheck"
	"github.com/surullabs/lint/govet"
)

func TestLint(t *testing.T) {
	linters := lint.Group(
		gofmt.Check{},
		govet.Shadow,
		golint.Check{},
		gosimple.Check{},
		gostaticcheck.Check{},
		dupl.Check{Threshold: 25},
	)

	linters = lint.Skip(linters,
		// Ignore all errors from unused.go
		lint.RegexpMatch(`unused\.go`),
		// Ignore duplicates we're okay with.
		dupl.SkipTwo, dupl.Skip("golint.go:1,12"),
	)
	if err := linters.Check("./..."); err != nil {
		t.Fatal(err)
	}
}

func Example() {
	linters := lint.Group(
		gofmt.Check{},         // Verify that all files are properly formatted
		govet.Shadow,          // go vet
		golint.Check{},        // golint
		gosimple.Check{},      // honnef.co/go/simple
		gostaticcheck.Check{}, // honnef.co/go/staticcheck
	)

	// Ignore some lint errors that we're not interested in. This ignores all errors from
	// the file unused.go. This is intended as an example of how to skip errors and not a
	// recommendation that you skip these kinds of errors.
	linters = lint.Skip(linters, lint.RegexpMatch(
		`unused\.go:4:2: a blank import`,
		`unused.go:7:7: don't use underscores in Go names`,
	))

	// Verify all files under this package recursively.
	if err := linters.Check("./..."); err != nil {
		// Record lint failures.
		// Use t.Fatal(err) when running in a test
		log.Fatal(err)
	}
	// Output:
}
