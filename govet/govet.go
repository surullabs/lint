package govet

import (
	"os/exec"
	"strings"

	"github.com/surullabs/lint/checkers"
)

// Check implements a lint.Checker for the govet command.
type Check struct {
	Args []string
}

// Shadow is a Checker that runs
// 	 go tool vet --all --shadow.
var Shadow = Check{Args: []string{"--all", "--shadow"}}

// Check runs go tool vet for pkgs.
func (c Check) Check(pkgs ...string) error {
	var errs []string
	for _, pkg := range pkgs {
		// Check files per package. If all files for all packages are
		// passed in as a glob, it causes incorrect reports as described
		// in TestGoVetMultiPackage_Issue7. Instead run go vet for each package.
		errs = append(errs, c.checkPackage(pkg)...)
	}
	return checkers.Error(errs...)
}

func (c Check) checkPackage(pkg string) []string {
	if strings.HasSuffix(pkg, "...") {
		return c.checkDir(pkg)
	}
	files, err := checkers.GoFiles(pkg)
	if err != nil {
		return []string{err.Error()}
	}
	return c.runVet(files)
}

func (c Check) runVet(paths []string) []string {
	if len(paths) == 0 {
		return nil
	}
	args := append([]string{"tool", "vet"}, append(c.Args, paths...)...)
	res, err := checkers.Exec(exec.Command("go", args...))
	if err == nil {
		return nil
	}
	switch res.Code {
	case 1:
		return strings.Split(strings.TrimSpace(res.Stderr), "\n")
	default:
		return []string{err.Error()}
	}
}

func (c Check) checkDir(pkg string) []string {
	p, err := checkers.Load(pkg)
	if err != nil {
		return []string{err.Error()}
	}
	// Loop through each package since we'd like to ignore directories which have
	// an _ prefix. In the future this allows us to skip vendor directories as well.
	var errs []string
	for _, pkg := range p.Pkgs {
		errs = append(errs, c.checkPackage(pkg)...)
	}
	return errs
}
