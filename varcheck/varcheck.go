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
	var args []string
	if c.ReportExported {
		args = append(args, "-e")
	}
	return checkers.Lint("varcheck", "github.com/opennota/check/cmd/varcheck", pkgs, args...)
}
