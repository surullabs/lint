package dupl

import (
	"testing"

	"github.com/surullabs/statictest/testutil"
)

func TestDupl(t *testing.T) {
	testutil.Test(t, "dupltest", []testutil.StaticCheckTest{
		{
			Checker:  Check{},
			Validate: testutil.NoError,
			Content: []byte(`package dupltest

import (
	"fmt"
)

func TestFunc() {
	fmt.Println("This is a file with no duplicates")
}
`),
		},
		{
			Checker:  Check{},
			Validate: testutil.MatchesRegexp("found 2 clones:"),
			Content: []byte(`package dupltest

import (
	"fmt"
)

func TestFunc() {
	fmt.Println("This is a duplicate string")
	fmt.Println("This is a duplicate string")
}

func TestFunc2() {
	fmt.Println("This is a duplicate string")
	fmt.Println("This is a duplicate string")
}
`),
		},
		{
			Checker: Check{},
			Content: []byte(`package dupltest

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
