package gostaticcheck

import (
	"github.com/surullabs/lint/checkers"
	_ "honnef.co/go/tools/staticcheck" // Ensure the staticcheck bin is downloaded.
)

// Check implements a gostaticcheck Checker (https://github.com/dominikh/go-staticcheck)
type Check struct {
	// Tags is a list of space separated build tags
	Tags string
}

// Check runs gostaticcheck for pkgs
func (c Check) Check(pkgs ...string) error {
	return checkers.Lint("staticcheck", "", "honnef.co/go/tools/cmd/staticcheck", pkgs, c.Args()...)
}

// Args returns command line arguments used for staticcheck
func (c Check) Args() []string {
	var args []string
	if c.Tags != "" {
		args = append(args, "-tags", c.Tags)
	}
	return args
}
