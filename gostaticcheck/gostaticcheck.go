package gostaticcheck

import "github.com/surullabs/statictest"

// Check implements a gostaticcheck Checker (https://github.com/dominikh/go-staticcheck)
type Check struct {
}

// Check runs gostaticcheck for pkg
func (Check) Check(pkg string) error {
	return statictest.Lint("staticcheck", "honnef.co/go/staticcheck/cmd/staticcheck", pkg)
}
