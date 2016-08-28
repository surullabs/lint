package gosimple

import (
	"github.com/surullabs/statictest/checkers"
	_ "honnef.co/go/simple" // Ensure the gosimple bin is downloaded.
)

// Check implements a gosimple Checker (https://github.com/dominikh/go-simple)
type Check struct {
}

// Check runs gosimple for pkg
func (Check) Check(pkgs ...string) error {
	return checkers.Lint("gosimple", "honnef.co/go/simple/cmd/gosimple", pkgs)
}
