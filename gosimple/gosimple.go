package gosimple

import (
	"github.com/surullabs/lint/checkers"
	_ "honnef.co/go/tools/simple" // Ensure the gosimple bin is downloaded.
)

// Check implements a gosimple Checker (https://github.com/dominikh/go-simple)
type Check struct {
	// Tags is a list of space separated build tags
	Tags string
}

// Check runs gosimple for pkg
func (c Check) Check(pkgs ...string) error {
	return checkers.Lint("gosimple", "", "honnef.co/go/tools/cmd/gosimple", pkgs, c.Args()...)
}

// Args returns command line arguments used for gosimple
func (c Check) Args() []string {
	var args []string
	if c.Tags != "" {
		args = append(args, "-tags", c.Tags)
	}
	return args
}
