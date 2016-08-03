package golint

import "github.com/surullabs/statictest/checkers"

// Check implements a golint Checker
type Check struct {
}

// Check runs golint for pkg
func (Check) Check(pkgs ...string) error {
	return checkers.LintAsFiles("golint", "github.com/golang/lint/golint", pkgs)
}
