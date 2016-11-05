package varcheck_test

import (
	"strings"
	"testing"

	"github.com/surullabs/lint"
	"github.com/surullabs/lint/varcheck"
	"github.com/surullabs/lint/testutil"
)

func TestGoVarcheck(t *testing.T) {
	testutil.Test(t, "varchecktest", []testutil.StaticCheckTest{
		{
			Checker: varcheck.Check{},
			Content: []byte(`package varchecktest
// TestFunc is a test function
func TestFunc() {
}
`),
			Validate: testutil.NoError,
		},
		{
			Checker: varcheck.Check{},
			Content: []byte(`package varchecktest
sfsff

func TestFunc() {
}
`),
			Validate: testutil.Contains("expected declaration, found 'IDENT' sfsff"),
		},
		{
			Checker: varcheck.Check{},
			Content: []byte(`package varchecktest
var unused bool
`),
			Validate: testutil.HasSuffix("unused"),
		},
		{
			Checker: lint.Skip(varcheck.Check{}, lint.StringSkipper{
				Strings: []string{
					"unused",
				},
				Matcher: strings.HasSuffix,
			}),
			Content: []byte(`package varchecktest

var unused bool
`),
			Validate: testutil.NoError,
		},
	},
	)
}
