package govet

import (
	"os/exec"
	"strings"

	"github.com/surullabs/statictest"
)

func Check(pkg string, args ...string) error {
	dir, err := statictest.PackageDir(pkg)
	if err != nil {
		return err
	}
	args = append([]string{"tool", "vet"}, append(args, dir)...)
	res, err := statictest.Exec(exec.Command("go", args...))
	if err == nil {
		return nil
	}
	switch res.Code {
	case 1:
		return &statictest.Error{Errors: strings.Split(strings.TrimSpace(res.Stderr), "\n")}
	default:
		return err
	}
}

func Checker(pkg string, args ...string) func(pkg string) error {
	return func(pkg string) error { return Check(pkg, args...) }
}
