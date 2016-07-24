package golint

import (
	"testing"

	"github.com/surullabs/statictest"
	"github.com/surullabs/statictest/testutil"
	"strings"
)

func TestGolint(t *testing.T) {
	testutil.Test(t, "golinttest", []testutil.StaticCheckTest{
		{
			Checker: Check{},
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
			Checker: Check{},
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
			Checker: Check{},
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
			Checker: statictest.Skip(Check{}, statictest.StringSkipper{
				Strings: []string{
					"exported function TestFunc should have comment or be unexported",
				},
				Matcher: strings.HasSuffix,
			}),
			Content: []byte(`package golinttest
import (
	"fmt"
)

func TestFunc() {
}
`),
			Validate: testutil.NoError,
		},
	},
	)
}
