package errcheck_test

import (
	"testing"

	"github.com/surullabs/lint/errcheck"
	"github.com/surullabs/lint/testutil"
)

func TestGoErrCheck(t *testing.T) {
	testutil.Test(t, "errchecktest", []testutil.StaticCheckTest{
		{
			Checker: errcheck.Check{},
			Content: []byte(`package errchecktest
// TestFunc is a test function
func TestFunc() {
}
`),
			Validate: testutil.NoError,
		},
		{
			Checker: errcheck.Check{},
			Content: []byte(`package errchecktest
sfsff

func TestFunc() {
}
`),
			Validate: testutil.Contains("expected declaration, found 'IDENT' sfsff"),
		},
		{
			Checker: errcheck.Check{},
			Content: []byte(`package errchecktest
import (
	"os"
)

func TestFunc() {
	f, _ := os.Open("somefile")
	f.Close()
}
`),
			Validate: testutil.HasSuffix("f.Close()"),
		},
		{
			Checker: errcheck.Check{},
			Content: []byte(`package errchecktest
import (
	"os"
)

func TestFunc() {
	f, _ := os.Open("somefile")
	f.Close()
}
`),
			Validate: testutil.SkippedErrors(`f\.Close`),
		},
		{
			Checker: errcheck.Check{Assert: true},
			Content: []byte(`package errchecktest

func TestFunc() {
	var i interface{} = 1
	_ = i.(int)
}
`),
			Validate: testutil.Contains("_ = i.(int)"),
		},
	},
	)
}

func TestGoErrCheckMultiFile(t *testing.T) {
	test := testutil.StaticCheckMultiFileTest{
		Contents: [][]byte{
			[]byte(`package errchecktest

func main() {
	f()
}
`),
			[]byte(`// +build tag

package errchecktest

import "errors"

func f() error {
	return errors.New("Returning error")
}
`),
		},
		Checker:  errcheck.Check{Tags: "tag"},
		Validate: testutil.HasSuffix("f()"),
	}

	if err := test.Test("errchecktest"); err != nil {
		t.Error("Check", err)
	}
}

func TestArgs(t *testing.T) {
	testutil.TestArgs(t, []testutil.ArgTest{
		{A: errcheck.Check{}, Expected: nil},
		{A: errcheck.Check{Blank: true}, Expected: []string{"-blank"}},
		{A: errcheck.Check{Assert: true}, Expected: []string{"-asserts"}},
		{A: errcheck.Check{Tags: "test"}, Expected: []string{"-tags", "test"}},
		{A: errcheck.Check{Blank: true, Assert: true}, Expected: []string{"-blank", "-asserts"}},
	})
}
