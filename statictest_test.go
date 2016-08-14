package statictest_test

import (
	"testing"

	"github.com/surullabs/statictest"
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
	)
	if err := statictest.Verify(basic, "."); err != nil {
		t.Fatal(err)
	}
}
