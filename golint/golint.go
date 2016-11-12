package golint

import "github.com/surullabs/lint/checkers"

// Check implements a golint Checker
type Check struct {
}

// Check implements lint.Checker for golint.
func (Check) Check(pkgs ...string) error {
	return checkers.Lint("golint", "", "github.com/golang/lint/golint", pkgs)
}
