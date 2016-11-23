package govet_test

import (
	"fmt"
	"testing"

	"strings"

	"path/filepath"

	"github.com/sridharv/fakegopath"
	"github.com/surullabs/lint/govet"
	"github.com/surullabs/lint/testutil"
)

func testVetError(err error) error {
	type errors interface {
		Errors() []string
	}
	skippable, ok := err.(errors)
	if !ok {
		return fmt.Errorf("unexpected type of error: %v", err)
	}
	errs := skippable.Errors()
	if len(errs) != 2 {
		return fmt.Errorf("expected 2 errors, got: %v", err)
	}
	if !strings.HasSuffix(errs[0], "unreachable code") {
		return err
	}
	if !strings.Contains(errs[1], "result of fmt.Sprintf call not used") {
		return err
	}
	return nil
}

// Using file samples from https://github.com/golang/go/issues/18018
const file1 = `package main

var foo error

func main() {
}`

const file2 = `package main

import (
        "errors"
        "fmt"
)

func main() {
        foo := errors.New("meomeo")
        fmt.Println(foo)
}`

const file3 = `package main

import "errors"

type foo error

func newFoo(msg string) foo {
        return foo(errors.New(msg))
}

func main() {
}`

func TestGoVetMultiPackage_Issue7(t *testing.T) {
	tmp, err := fakegopath.NewTemporaryWithFiles("multipkg", []fakegopath.SourceFile{
		{Content: []byte(file1), Dest: filepath.Join("root", "package1", "main.go")},
		{Content: []byte(file2), Dest: filepath.Join("root", "package2", "main.go")},
		{Content: []byte(file3), Dest: filepath.Join("root", "package3", "main.go")},
	})
	if err != nil {
		t.Fatalf("failed to create temporary go path: %v", err)
	}
	defer tmp.Reset()

	if err := govet.Shadow.Check("root/..."); err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if err := govet.Shadow.Check("root/package1", "root/package2", "root/package3"); err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
}

func TestGoVet(t *testing.T) {
	testutil.Test(t, "govettest", []testutil.StaticCheckTest{
		{
			Checker: govet.Check{},
			Content: []byte(`package govettest
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
			Checker: govet.Check{},
			Content: []byte(`package govettest

import (
	"fmt"
)
sfsff

func TestFunc() {
	fmt.Println("undocumented")
}
`),
			Validate: testutil.Contains("expected declaration, found 'IDENT' sfsff"),
		},
		{
			Checker: govet.Check{},
			Content: []byte(`package gofmttest

import (
	"fmt"
)

func TestFunc() {
    a := "test"
    b := a
    fmt.Sprintf("test")
    return
    fmt.Println("This is a poorly formatted file", b)
}
`),
			Validate: testVetError,
		},
		{
			Checker: govet.Shadow,
			Content: []byte(`package gofmttest

import (
	"fmt"
)

func TestFunc() (err error) {
    err = fmt.Println("another")
    if err != nil {
    	err := fmt.Errorf("some error: %v", err)
    }
    return err
}
`),
			Validate: testutil.MatchesRegexp(`declaration of "?err"? shadows declaration at`),
		},
	})

}
