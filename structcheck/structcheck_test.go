package structcheck_test

import (
	"strings"
	"testing"

	"github.com/surullabs/lint"
	"github.com/surullabs/lint/structcheck"
	"github.com/surullabs/lint/testutil"
)

func TestStructcheck(t *testing.T) {
	testutil.Test(t, "structchecktest", []testutil.StaticCheckTest{
		{
			Checker: structcheck.Check{},
			Content: []byte(`package structchecktest
// TestFunc is a test function
func TestFunc() {
}
`),
			Validate: testutil.NoError,
		},
		{
			Checker: structcheck.Check{},
			Content: []byte(`package structchecktest
sfsff

func TestFunc() {
}
`),
			Validate: testutil.Contains("expected declaration, found 'IDENT' sfsff"),
		},
		{
			Checker: structcheck.Check{},
			Content: []byte(`package structchecktest
type s struct {
	b bool
}
`),
			Validate: testutil.HasSuffix("structchecktest.s.b"),
		},
		{
			Checker: lint.Skip(structcheck.Check{}, lint.StringSkipper{
				Strings: []string{
					"structchecktest.s.b",
				},
				Matcher: strings.HasSuffix,
			}),
			Content: []byte(`package structchecktest
type s struct {
	b bool
}
`),
			Validate: testutil.NoError,
		},
	},
	)
}
