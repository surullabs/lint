package statictest_test

import (
	"testing"

	"github.com/surullabs/statictest"
	"github.com/surullabs/statictest/dupl"
	"github.com/surullabs/statictest/gofmt"
	"github.com/surullabs/statictest/golint"
	"github.com/surullabs/statictest/gosimple"
	"github.com/surullabs/statictest/gostaticcheck"
	"github.com/surullabs/statictest/govet"
)

func TestStaticChecks(t *testing.T) {
	basic := statictest.Group(
		gofmt.Check{},
		govet.Shadow,
		golint.Check{},
		gosimple.Check{},
		gostaticcheck.Check{},
		dupl.Check{Threshold: 25},
	)
	// Ignore duplicates we're okay with.
	skipped := statictest.Skip(basic, dupl.SkipTwo, dupl.Skip("golint.go:1,12"))
	if err := skipped.Check("./..."); err != nil {
		t.Fatal(err)
	}
}
