package dupl_test

import (
	"testing"

	"github.com/surullabs/lint/dupl"
	"github.com/surullabs/lint/testutil"
)

func TestDupl(t *testing.T) {
	testutil.Test(t, "dupltest", []testutil.StaticCheckTest{
		{
			Checker:  dupl.Check{},
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
			Checker:  dupl.Check{},
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
			Checker: dupl.Check{},
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

func TestSkipTwo(t *testing.T) {
	testutil.TestSkips(t, []testutil.SkipTest{
		{S: dupl.SkipTwo, Line: "some line", Skip: false},
		{S: dupl.SkipTwo, Line: "found 2 clones: here", Skip: true},
		{S: dupl.SkipTwo, Line: "checker: found 2 clones: here", Skip: false},
		{S: dupl.SkipTwo, Line: "dupl.Check: found 2 clones: here", Skip: true},
	})
}

func TestSkip(t *testing.T) {
	testutil.TestSkips(t, []testutil.SkipTest{
		{S: dupl.Skip("lint.go:1,12"), Line: "some line", Skip: false},
		{S: dupl.Skip("lint.go:1,12"), Line: "lint.go:1,12", Skip: false},
		{S: dupl.Skip("lint.go:1,12"), Line: "dupl.Check: found 2 clones: here\nlint.go:1,12", Skip: true},
	})
}
