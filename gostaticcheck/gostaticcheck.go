package gostaticcheck

import (
	"github.com/surullabs/lint/checkers"
	_ "honnef.co/go/staticcheck" // Ensure the staticcheck bin is downloaded.
)

// Check implements a gostaticcheck Checker (https://github.com/dominikh/go-staticcheck)
type Check struct {
}

// Check runs gostaticcheck for pkgs
func (Check) Check(pkgs ...string) error {
	return checkers.Lint("staticcheck", "honnef.co/go/staticcheck/cmd/staticcheck", pkgs)
}
