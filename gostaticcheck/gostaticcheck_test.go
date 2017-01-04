package gostaticcheck_test

import (
	"testing"

	"github.com/surullabs/lint/gostaticcheck"
	"github.com/surullabs/lint/testutil"
)

func TestGostaticcheck(t *testing.T) {
	testutil.Test(t, "gostaticchecktest", []testutil.StaticCheckTest{
		{
			Checker: gostaticcheck.Check{},
			Content: []byte(`package gostaticchecktest
import (
	"fmt"
)

// TestFunc is a test function
func TestFunc() {
	fmt.Println("This is a properly formatted file")
}
`),
			Validate: testutil.NoError,
		},
		{
			Checker: gostaticcheck.Check{},
			Content: []byte(`package gostaticchecktest

import (
	"fmt"
)
sfsff

func TestFunc() {
	fmt.Println("undocumented")
}
`),
			Validate: testutil.Contains("expected declaration, found 'IDENT' sfsff"),
		},
		{
			Checker: gostaticcheck.Check{},
			Content: []byte(`package gostaticchecktest
import (
	"regexp"
)

func TestFunc() {
	regexp.Compile("foo(")
}
`),
			Validate: testutil.Contains(
				" error parsing regexp: missing closing ): `foo(`"),
		},
	},
	)
}

func TestGostaticcheckMultiFile(t *testing.T) {
	test := testutil.StaticCheckMultiFileTest{
		Contents: [][]byte{
			[]byte(`package gostaticchecktest

func main() {
	f()
}
`),
			[]byte(`// +build tag

package gostaticchecktest

func f() {
	b := true
	if !!b {
	}
}
`),
		},
		Checker:  gostaticcheck.Check{Tags: "tag"},
		Validate: testutil.Contains(" negating a boolean twice has no effect"),
	}

	if err := test.Test("gostaticchecktest"); err != nil {
		t.Error("Check", err)
	}
}

func TestArgs(t *testing.T) {
	testutil.TestArgs(t, []testutil.ArgTest{
		{A: gostaticcheck.Check{}, Expected: nil},
		{A: gostaticcheck.Check{Tags: "test"}, Expected: []string{"-tags", "test"}},
	})
}
