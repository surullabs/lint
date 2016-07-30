package golint

import "github.com/surullabs/statictest/checkers"

// Check implements a golint Checker
type Check struct {
}

// Check runs golint for pkg
func (Check) Check(pkg string) error {
	return checkers.Lint("golint", "github.com/golang/lint/golint", pkg)
}
