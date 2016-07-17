package statictest_test

import (
	"testing"

	"github.com/surullabs/statictest/gofmt"
	"github.com/surullabs/statictest/govet"
)

func TestStaticChecks(t *testing.T) {
	type checkFn func(string) error
	for _, check := range []checkFn{
		gofmt.Check, govet.Checker(),
	} {
		if err := check("."); err != nil {
			t.Error(err)
		}
	}
}
