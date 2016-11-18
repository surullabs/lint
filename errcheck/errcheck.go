// Package errcheck provides lint integration for the errcheck linter
package errcheck

import "github.com/surullabs/lint/checkers"

// Check runs the errcheck linter (https://github.com/kisielk/errcheck)
type Check struct {
	// Blank enables checking for assignments to the blank identifier
	Blank bool
	// Asserts enables checking for ignored type assertion results
	Assert bool
	// Tags is a list of space separated build tags
	Tags string
}

// Check runs errcheck and returns any errors found.
func (c Check) Check(pkgs ...string) error {
	return checkers.Lint("errcheck", "", "github.com/kisielk/errcheck", pkgs, c.Args()...)
}

// Args returns command line arguments used for errcheck
func (c Check) Args() []string {
	var args []string
	if c.Blank {
		args = append(args, "-blank")
	}
	if c.Assert {
		args = append(args, "-asserts")
	}
	if c.Tags != "" {
		args = append(args, "-tags", c.Tags)
	}
	return args
}
