package golint_test

import (
	"testing"

	"github.com/surullabs/lint/golint"
	"github.com/surullabs/lint/testutil"
)

func TestGolint(t *testing.T) {
	testutil.Test(t, "golinttest", []testutil.StaticCheckTest{
		{
			Checker: golint.Check{},
			Content: []byte(`package golinttest
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
			Checker: golint.Check{},
			Content: []byte(`package golinttest

import (
	"fmt"
)
sfsff

func TestFunc() {
	fmt.Println("undocumented")
}
`),
			Validate: testutil.HasSuffix("expected declaration, found 'IDENT' sfsff"),
		},
		{
			Checker: golint.Check{},
			Content: []byte(`package golinttest
import (
	"fmt"
)

func TestFunc() {
	fmt.Println("This is a properly formatted file")
}
`),
			Validate: testutil.HasSuffix(
				"file.go:6:1: exported function TestFunc should have comment or be unexported"),
		},
		{
			Checker: golint.Check{},
			Content: []byte(`package golinttest
import (
	"fmt"
)

func TestFunc() {
}
`),
			Validate: testutil.SkippedErrors(
				`exported function TestFunc should have comment or be unexported`),
		},
	},
	)
}
