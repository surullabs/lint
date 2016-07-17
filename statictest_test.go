package statictest_test

import (
	"testing"

	"github.com/surullabs/statictest"
	"github.com/surullabs/statictest/gofmt"
	"github.com/surullabs/statictest/govet"
)

func TestStaticChecks(t *testing.T) {
	statictest.Check(t, ".", gofmt.Check, govet.Checker())
}
