package gometalinter_test

import (
	"testing"

	"log"

	"github.com/surullabs/lint/gometalinter"
	"github.com/surullabs/lint/testutil"
)

func TestMetaLinter(t *testing.T) {
	testutil.Test(t, "gometalintertest", []testutil.StaticCheckTest{
		{
			Checker: gometalinter.Check{},
			Content: []byte(`package gometalintertest
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
			Checker: gometalinter.Check{},
			Content: []byte(`package gometalintertest

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
			Checker: gometalinter.Check{},
			Content: []byte(`package gometalintertest
import (
	"fmt"
)

func TestFunc() {
	fmt.Println("This is a properly formatted file")
}
`),
			Validate: testutil.HasSuffix(
				"file.go:6:1:warning: exported function TestFunc should have comment or be unexported (golint)"),
		},
		{
			Checker: gometalinter.Check{},
			Content: []byte(`package gometalintertest
func TestFunc() {
}
`),
			Validate: testutil.SkippedErrors(
				`exported function TestFunc should have comment or be unexported`),
		},
	},
	)
}

func Example() {
	metalinter := gometalinter.Check{
		Args: []string{
			// These are not recommendations for linters to disable.
			"--disable=gocyclo",
			"--disable=gas",
			"--deadline=20s",
		},
	}
	if err := metalinter.Check("./..."); err != nil {
		log.Fatal(err)
	}

	// Output:
}
