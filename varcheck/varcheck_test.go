package varcheck_test

import (
	"testing"

	"github.com/surullabs/lint/testutil"
	"github.com/surullabs/lint/varcheck"
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
			Checker: varcheck.Check{},
			Content: []byte(`package varchecktest

var unused bool
`),
			Validate: testutil.SkippedErrors("unused"),
		},
	},
	)
}

func TestArgs(t *testing.T) {
	testutil.TestArgs(t, []testutil.ArgTest{
		{A: varcheck.Check{}, Expected: nil},
		{A: varcheck.Check{ReportExported: true}, Expected: []string{"-e"}},
	})
}
