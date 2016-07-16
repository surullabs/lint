package statictest_test

import (
	"github.com/surullabs/statictest/gofmt"
	"testing"
)

func TestStaticChecks(t *testing.T) {
	type checkFn func(string) error
	for _, check := range []checkFn{
		gofmt.Check,
	} {
		if err := check("."); err != nil {
			t.Error(err)
		}
	}
}
