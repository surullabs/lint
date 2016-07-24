package gofmt

import (
	"testing"

	"regexp"

	"github.com/surullabs/statictest/testutil"
)

var unformattedRE = regexp.MustCompile("/var/[^ ]*gofmt(test)?[0-9]+.*\n")

func testUnformatted(err error) error {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return testutil.Diff(expectedUnformatted, unformattedRE.ReplaceAllString(errStr, "GOFMT_TMP_FOLDER\n"))
}

func TestGoFmt(t *testing.T) {
	testutil.Test(t, "gofmttest", []testutil.StaticCheckTest{
		{
			Checker: Check{},
			Content: []byte(`package gofmttest

import (
	"fmt"
)

func TestFunc() {
	fmt.Println("This is a properly formatted file")
}
`),
			Validate: testutil.NoError,
		},
		{
			Checker: Check{},
			Content: []byte(`package gofmttest

import (
	"fmt"
)

func TestFunc() {
  fmt.Println("This is a poorly formatted file")
}

type A struct {
    A string
    Long string
}
`),
			Validate: testUnformatted,
		},
		{
			Checker: Check{},
			Content: []byte(`package gofmttest

blah`),
			Validate: testutil.HasSuffix("expected declaration, found 'IDENT' blah\n"),
		},
	})

}

const expectedUnformatted = `File not formatted: diff GOFMT_TMP_FOLDER
--- GOFMT_TMP_FOLDER
+++ GOFMT_TMP_FOLDER
@@ -5,10 +5,10 @@
 )
 
 func TestFunc() {
-  fmt.Println("This is a poorly formatted file")
+	fmt.Println("This is a poorly formatted file")
 }
 
 type A struct {
-    A string
-    Long string
+	A    string
+	Long string
 }`
