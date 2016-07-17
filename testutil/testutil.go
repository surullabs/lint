package testutil

import (
	"fmt"

	"path/filepath"

	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/sridharv/fakegopath"
)

type StaticCheckTest struct {
	File     string
	Args     interface{}
	Validate func(err error) error
}

func (s StaticCheckTest) Test(pkg string, check func(string, interface{}) error) error {
	tmp, err := fakegopath.NewTemporaryWithFiles(pkg, []fakegopath.SourceFile{
		{Src: s.File, Dest: filepath.Join(pkg, "file.go")},
	})
	if err != nil {
		return fmt.Errorf("failed to create temporary go path: %v", err)
	}
	defer tmp.Reset()
	return s.Validate(check(pkg, s.Args))
}

type Errorer interface {
	Error(args ...interface{})
}

func Test(t Errorer, pkg string, check func(string, interface{}) error, tests []StaticCheckTest) {
	for _, test := range tests {
		if err := test.Test(pkg, check); err != nil {
			t.Error(test.File, test.Args, err)
		}
	}
}

func NoError(err error) error { return err }

func HasSuffix(suffix string) func(err error) error {
	return func(err error) error {
		if err == nil {
			return fmt.Errorf("no error found when expecting error with suffix %s", suffix)
		}
		if !strings.HasSuffix(err.Error(), suffix) {
			return err
		}
		return nil
	}
}

func Diff(expected, actual string) error {
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
	return fmt.Errorf(text)
}
