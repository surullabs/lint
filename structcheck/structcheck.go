// Package structcheck provides lint integration for the structcheck linter
package structcheck

import "github.com/surullabs/lint/checkers"

// Check runs the structcheck linter (https://github.com/opennota/check)
type Check struct {
	// ReportExported reports exported fields that are unused
	ReportExported bool
	// OnlyCountAssignments ensures only assignments are counted
	OnlyCountAssignments bool
	// IncludeTests loads test files
	IncludeTests bool
}

// Check runs structcheck and returns any errors found.
func (c Check) Check(pkgs ...string) error {
	return checkers.Lint("structcheck",
		"github.com/opennota/check",
		"github.com/opennota/check/cmd/structcheck", pkgs, c.Args()...)
}

// Args returns command line flags for structcheck
func (c Check) Args() []string {
	var args []string
	if c.ReportExported {
		args = append(args, "-e")
	}
	if c.OnlyCountAssignments {
		args = append(args, "-a")
	}
	if c.IncludeTests {
		args = append(args, "-t")
	}
	return args
}
