package govet

import (
	"os/exec"
	"strings"

	"github.com/surullabs/statictest"
)

// Check implements a statictest.Checker for the govet command.
type Check struct {
	Args []string
}

func (c Check) Check(pkg string) error {
	dir, err := statictest.PackageDir(pkg)
	if err != nil {
		return err
	}
	args := append([]string{"tool", "vet"}, append(c.Args, dir)...)
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
