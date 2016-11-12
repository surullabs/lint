package structcheck_test

import (
	"testing"

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
			Checker: structcheck.Check{},
			Content: []byte(`package structchecktest
type s struct {
	b bool
}
`),
			Validate: testutil.SkippedErrors(`structchecktest\.s\.b`),
		},
	},
	)
}

func TestArgs(t *testing.T) {
	testutil.TestArgs(t, []testutil.ArgTest{
		{A: structcheck.Check{}, Expected: nil},
		{A: structcheck.Check{IncludeTests: true}, Expected: []string{"-t"}},
		{A: structcheck.Check{OnlyCountAssignments: true}, Expected: []string{"-a"}},
		{A: structcheck.Check{ReportExported: true}, Expected: []string{"-e"}},
		{A: structcheck.Check{IncludeTests: true, ReportExported: true}, Expected: []string{"-e", "-t"}},
	})
}
