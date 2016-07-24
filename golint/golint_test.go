package golint

import (
	"testing"

	"github.com/surullabs/statictest/testutil"
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
			Validate: testutil.NoError,
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
	},
	)
}
