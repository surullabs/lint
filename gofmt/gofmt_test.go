package gofmt

import (
	"testing"

	"fmt"

	"regexp"

	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/sridharv/fakegopath"
)

type goFmtTest struct {
	file    string
	checkFn func(err error) error
}

func (g goFmtTest) test() error {
	tmp, err := fakegopath.NewTemporaryWithFiles("gofmt", []fakegopath.SourceFile{
		{Src: g.file, Dest: "gofmttest/file.go"},
	})
	if err != nil {
		return fmt.Errorf("failed to create temporary go path: %v", err)
	}
	defer tmp.Reset()
	return g.checkFn(Check("gofmttest"))
}

func testClean(err error) error { return err }

var unformattedRE = regexp.MustCompile("/var/[^ ]*gofmt[0-9]+.*\n")

func diff(expected, actual string) error {
	if expected == actual {
		return nil
	}
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(actual),
		FromFile: "Golden",
		ToFile:   "Actual",
		Context:  3,
	}
	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		text = fmt.Sprintf("diff error: %v", err)
	}
	return fmt.Errorf("(golden) %s != actual\n%s", expected, text)
}

func testUnformatted(err error) error {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return diff(expectedUnformatted, unformattedRE.ReplaceAllString(errStr, "GOFMT_TMP_FOLDER\n"))
}

func testCompilerError(err error) error {
	if strings.HasSuffix(err.Error(), "src/gofmttest/file.go:3:1: expected declaration, found 'IDENT' blah\n") {
		return nil
	}
	return err
}

func TestGoFmt(t *testing.T) {
	tests := []goFmtTest{
		{"testdata/clean.go.src", testClean},
		{"testdata/unformatted.go.src", testUnformatted},
		{"testdata/compilererror.go.src", testCompilerError},
	}
	for _, test := range tests {
		if err := test.test(); err != nil {
			t.Errorf("testing %s: %v", test.file, err)
		}
	}
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
