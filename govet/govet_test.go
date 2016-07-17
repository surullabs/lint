package govet

import (
	"fmt"
	"testing"

	"github.com/sridharv/fakegopath"
	"strings"
	"github.com/surullabs/statictest"
)

type goVetTest struct {
	file    string
	args    []string
	checkFn func(err error) error
}

func (g goVetTest) test() error {
	tmp, err := fakegopath.NewTemporaryWithFiles("govet", []fakegopath.SourceFile{
		{Src: g.file, Dest: "govettest/file.go"},
	})
	if err != nil {
		return fmt.Errorf("failed to create temporary go path: %v", err)
	}
	defer tmp.Reset()
	return g.checkFn(Check("govettest", g.args...))
}

func testClean(err error) error { return err }

func testCompilerError(err error) error {
	if !strings.HasSuffix(err.Error(), "expected declaration, found 'IDENT' sfsff") {
		return err
	}
	return nil
}

func testVetError(err error) error {
	skippable, ok := err.(*statictest.Error)
	if !ok {
		return fmt.Errorf("unexpected type of error: %v", err)
	}
	if len(skippable.Errors) != 2 {
		return fmt.Errorf("expected 2 errors, got: %v", err)
	}
	if !strings.HasSuffix(skippable.Errors[0], "unreachable code") {
		return err
	}
	if !strings.HasSuffix(skippable.Errors[1], "result of fmt.Sprintf call not used") {
		return err
	}
	return nil
}

func TestGoVet(t *testing.T) {
	tests := []goVetTest{
		{"testdata/clean.go.src", []string{}, testClean},
		{"testdata/compilererror.go.src", []string{}, testCompilerError},
		{"testdata/veterrors.go.src", []string{}, testVetError},
	}
	for _, test := range tests {
		if err := test.test(); err != nil {
			t.Errorf("testing %s: %v", test.file, err)
		}
	}
}

//var unformattedRE = regexp.MustCompile("/var/[^ ]*gofmt[0-9]+.*\n")
//
//func diff(expected, actual string) error {
//	if expected == actual {
//		return nil
//	}
//	diff := difflib.UnifiedDiff{
//		A:        difflib.SplitLines(expected),
//		B:        difflib.SplitLines(actual),
//		FromFile: "Golden",
//		ToFile:   "Actual",
//		Context:  3,
//	}
//	text, err := difflib.GetUnifiedDiffString(diff)
//	if err != nil {
//		text = fmt.Sprintf("diff error: %v", err)
//	}
//	return fmt.Errorf("(golden) %s != actual\n%s", expected, text)
//}
//
//func testUnformatted(err error) error {
//	errStr := ""
//	if err != nil {
//		errStr = err.Error()
//	}
//	return diff(expectedUnformatted, unformattedRE.ReplaceAllString(errStr, "GOFMT_TMP_FOLDER\n"))
//}
//
//func testCompilerError(err error) error {
//	if strings.HasSuffix(err.Error(), "src/gofmttest/file.go:3:1: expected declaration, found 'IDENT' blah\n") {
//		return nil
//	}
//	return err
//}
