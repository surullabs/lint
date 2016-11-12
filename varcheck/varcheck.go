// Package varcheck provides lint integration for the varcheck linter
package varcheck

import "github.com/surullabs/lint/checkers"

// Check runs the varcheck linter (https://github.com/opennota/check)
type Check struct {
	// ReportExported reports exported variables that are unused
	ReportExported bool
}

// Check runs varcheck and returns any errors found.
func (c Check) Check(pkgs ...string) error {
	if _, err := checkers.InstallMissing("varcheck", "github.com/opennota/check", "github.com/opennota/check/cmd/varcheck"); err != nil {
		return err
	}
	return checkers.Lint("varcheck",
		"github.com/opennota/check",
		"github.com/opennota/check/cmd/varcheck", pkgs, c.Args()...)
}

// Args returns all args passed to varcheck
func (c Check) Args() []string {
	var args []string
	if c.ReportExported {
		args = append(args, "-e")
	}
	return args
}
