package golint

import "github.com/surullabs/statictest/checkers"

// Check implements a golint Checker
type Check struct {
}

// Check implements statictest.Checker for golint.
func (Check) Check(pkgs ...string) error {
	return checkers.Lint("golint", "github.com/golang/lint/golint", pkgs)
}
