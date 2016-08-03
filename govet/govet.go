package govet

import (
	"os/exec"
	"strings"

	"github.com/surullabs/statictest"
	"github.com/surullabs/statictest/checkers"
)

// Check implements a statictest.Checker for the govet command.
type Check struct {
	Args []string
}

var Shadow = Check{Args: []string{"--all", "--shadow"}}

func (c Check) Check(pkgs ...string) error {
	files, err := checkers.GoFiles(pkgs...)
	if err != nil {
		return err
	}
	args := append([]string{"tool", "vet"}, append(c.Args, files...)...)
	res, err := checkers.Exec(exec.Command("go", args...))
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
