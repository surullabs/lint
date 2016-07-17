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
	check := func(pkg string, _ interface{}) error { return Check(pkg) }
	testutil.Test(t, "gofmttest", check, []testutil.StaticCheckTest{
		{"testdata/clean.go.src", nil, testutil.NoError},
		{"testdata/unformatted.go.src", nil, testUnformatted},
		{"testdata/compilererror.go.src", nil, testutil.HasSuffix("expected declaration, found 'IDENT' blah\n")},
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
