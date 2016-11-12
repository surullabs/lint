package govet_test

import (
	"fmt"
	"testing"

	"strings"

	"github.com/surullabs/lint/testutil"
	"github.com/surullabs/lint/govet"
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
