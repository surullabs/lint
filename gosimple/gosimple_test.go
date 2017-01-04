package gosimple_test

import (
	"testing"

	"github.com/surullabs/lint/gosimple"
	"github.com/surullabs/lint/testutil"
)

func TestGosimple(t *testing.T) {
	testutil.Test(t, "gosimpletest", []testutil.StaticCheckTest{
		{
			Checker: gosimple.Check{},
			Content: []byte(`package gosimpletest
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
			Checker: gosimple.Check{},
			Content: []byte(`package gosimpletest

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
			Checker: gosimple.Check{},
			Content: []byte(`package gosimpletest

func TestFunc() {
	for _ = range []string{"a", "b"} {
	}
}
`),
			Validate: testutil.Contains(
				"should omit values from range; this loop is equivalent to `for range ...`"),
		},
	},
	)
}

func TestGosimpleMultiFile(t *testing.T) {
	test := testutil.StaticCheckMultiFileTest{
		Contents: [][]byte{
			[]byte(`package gosimpletest

func main() {
	f()
}
`),
			[]byte(`// +build tag

package gosimpletest

func f() {
	b := true
	if b == true {
	}
}
`),
		},
		Checker:  gosimple.Check{Tags: "tag"},
		Validate: testutil.Contains(" should omit comparison to bool constant"),
	}

	if err := test.Test("gosimpletest"); err != nil {
		t.Error("Check", err)
	}
}

func TestArgs(t *testing.T) {
	testutil.TestArgs(t, []testutil.ArgTest{
		{A: gosimple.Check{}, Expected: nil},
		{A: gosimple.Check{Tags: "test"}, Expected: []string{"-tags", "test"}},
	})
}
