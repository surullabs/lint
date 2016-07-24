package gofmt

import (
	"testing"

	"github.com/surullabs/statictest/testutil"
)

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
			Validate: testutil.MatchesRegexp("^File not formatted.*gofmttest/file.go"),
		},
		{
			Checker: Check{},
			Content: []byte(`package gofmttest

		blah`),
			Validate: testutil.Contains("expected declaration, found 'IDENT' blah"),
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
