package gosimple

import (
	"testing"

	"github.com/surullabs/statictest/testutil"
)

func TestGosimple(t *testing.T) {
	testutil.Test(t, "gosimpletest", []testutil.StaticCheckTest{
		{
			Checker: Check{},
			Content: []byte(`package gosimpletest
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
			Content: []byte(`package gosimpletest

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
			Content: []byte(`package gosimpletest
import (
	"fmt"
)

func TestFunc() {
	for _ = range []string{"a", "b"} {
	}
}
`),
			Validate: testutil.HasSuffix(
				"should omit values from range; this loop is equivalent to `for range ...`"),
		},
	},
	)
}
