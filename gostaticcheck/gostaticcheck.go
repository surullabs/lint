package gostaticcheck

import "github.com/surullabs/statictest/checkers"

// Check implements a gostaticcheck Checker (https://github.com/dominikh/go-staticcheck)
type Check struct {
}

// Check runs gostaticcheck for pkgs
func (Check) Check(pkgs ...string) error {
	return checkers.LintAsFiles("staticcheck", "honnef.co/go/staticcheck/cmd/staticcheck", pkgs)
}
