package testutil

import (
	"fmt"

	"path/filepath"

	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/sridharv/fakegopath"
	"github.com/surullabs/statictest"
	"regexp"
)

type StaticCheckTest struct {
	File     string
	Content  []byte
	Checker  statictest.Checker
	Validate func(err error) error
}

func (s StaticCheckTest) Test(dir string) error {
	tmp, err := fakegopath.NewTemporaryWithFiles(dir, []fakegopath.SourceFile{
		{Src: s.File, Content: s.Content, Dest: filepath.Join(dir, "file.go")},
	})
	if err != nil {
		return fmt.Errorf("failed to create temporary go path: %v", err)
	}
	defer tmp.Reset()
	return s.Validate(s.Checker.Check(dir))
}

type Errorer interface {
	Error(args ...interface{})
}

func Test(t Errorer, pkg string, tests []StaticCheckTest) {
	for i, test := range tests {
		if err := test.Test(pkg); err != nil {
			t.Error("Check", i, err)
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

func MatchesRegexp(re string) func (err error) error {
	return func (err error) error {
		if err == nil {
			return fmt.Errorf("no error found when expecting error matching RE %s", re)
		}
		if matches, matchErr := regexp.MatchString(re, err.Error()); matchErr != nil {
			return matchErr
		} else if !matches {
			return fmt.Errorf("error %v does not match re %s", err, re)
		}
		return nil
	}
}

func Contains(str string) func(err error) error {
	return func(err error) error {
		if err == nil {
			return fmt.Errorf("no error found when expecting error containing %s", str)
		}
		if !strings.Contains(err.Error(), str) {
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
