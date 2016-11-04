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

func TestStaticChecks(t *testing.T) {
	basic := lint.Group(
		gofmt.Check{},
		govet.Shadow,
		golint.Check{},
		gosimple.Check{},
		gostaticcheck.Check{},
		dupl.Check{Threshold: 25},
	)
	// Ignore duplicates we're okay with.
	skipped := lint.Skip(basic, dupl.SkipTwo, dupl.Skip("golint.go:1,12"))
	if err := skipped.Check("./..."); err != nil {
		t.Fatal(err)
	}
}
