package aligncheck_test

import (
	"strings"
	"testing"

	"github.com/surullabs/lint"
	"github.com/surullabs/lint/aligncheck"
	"github.com/surullabs/lint/testutil"
)

func TestAligncheck(t *testing.T) {
	testutil.Test(t, "alignchecktest", []testutil.StaticCheckTest{
		{
			Checker: aligncheck.Check{},
			Content: []byte(`package alignchecktest
// TestFunc is a test function
func TestFunc() {
}
`),
			Validate: testutil.NoError,
		},
		{
			Checker: aligncheck.Check{},
			Content: []byte(`package alignchecktest
sfsff

func TestFunc() {
}
`),
			Validate: testutil.Contains("expected declaration, found 'IDENT' sfsff"),
		},
		{
			Checker: aligncheck.Check{},
			Content: []byte(`package alignchecktest

type s struct {
	b bool
	a string
	c int32
}
`),
			Validate: testutil.HasSuffix("struct s could have size 24 (currently 32)"),
		},
		{
			Checker: lint.Skip(aligncheck.Check{}, lint.StringSkipper{
				Strings: []string{
					"struct s could have size 24 (currently 32)",
				},
				Matcher: strings.HasSuffix,
			}),
			Content: []byte(`package alignchecktest

type s struct {
	b bool
	a string
	c int32
}
`),
			Validate: testutil.NoError,
		},
	},
	)
}
