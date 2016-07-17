package govet

import (
	"fmt"
	"testing"

	"strings"

	"github.com/surullabs/statictest"
	"github.com/surullabs/statictest/testutil"
)

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
	check := func(pkg string, args interface{}) error { return Check(pkg, args.([]string)...) }
	testutil.Test(t, "govettest", check, []testutil.StaticCheckTest{
		{"testdata/clean.go.src", []string{}, testutil.NoError},
		{"testdata/compilererror.go.src", []string{}, testutil.HasSuffix("expected declaration, found 'IDENT' sfsff")},
		{"testdata/veterrors.go.src", []string{}, testVetError},
	})
}
