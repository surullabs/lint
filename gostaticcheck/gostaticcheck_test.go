package gostaticcheck

import (
	"testing"

	"github.com/surullabs/statictest/testutil"
)

func TestGostaticcheck(t *testing.T) {
	testutil.Test(t, "gostaticchecktest", []testutil.StaticCheckTest{
		{
			Checker: Check{},
			Content: []byte(`package gostaticchecktest
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
			Content: []byte(`package gostaticchecktest

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
			Content: []byte(`package gostaticchecktest
import (
	"fmt"
	"regexp"
)

func TestFunc() {
	regexp.Compile("foo(")
}
`),
			Validate: testutil.HasSuffix(
				" error parsing regexp: missing closing ): `foo(`"),
		},
	},
	)
}
