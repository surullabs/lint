package errcheck_test

import (
	"strings"
	"testing"

	"github.com/surullabs/lint"
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
			Checker: lint.Skip(errcheck.Check{}, lint.StringSkipper{
				Strings: []string{
					"f.Close()",
				},
				Matcher: strings.HasSuffix,
			}),
			Content: []byte(`package errchecktest
import (
	"os"
)

func TestFunc() {
	f, _ := os.Open("somefile")
	f.Close()
}
`),
			Validate: testutil.NoError,
		},
	},
	)
}
