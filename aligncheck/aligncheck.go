// Package aligncheck provides lint integration for the aligncheck linter
package aligncheck

import "github.com/surullabs/lint/checkers"

// Check runs the aligncheck linter (https://github.com/opennota/check)
type Check struct {
}

// Check runs aligncheck and returns any errors found.
func (c Check) Check(pkgs ...string) error {
	return checkers.Lint("aligncheck",
		"github.com/opennota/check",
		"github.com/opennota/check/cmd/aligncheck", pkgs)
}
