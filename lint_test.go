package lint_test

import (
	"testing"

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
	// Ignore duplicates we're okay with.
	linters = lint.Skip(linters, dupl.SkipTwo, dupl.Skip("golint.go:1,12"))
	if err := linters.Check("./..."); err != nil {
		t.Fatal(err)
	}
}
