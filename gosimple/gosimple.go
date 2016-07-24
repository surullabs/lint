package gosimple

import "github.com/surullabs/statictest"

// Check implements a gosimple Checker (https://github.com/dominikh/go-simple)
type Check struct {
}

// Check runs gosimple for pkg
func (Check) Check(pkg string) error {
	return statictest.Lint("gosimple", "honnef.co/go/simple/cmd/gosimple", pkg)
}
