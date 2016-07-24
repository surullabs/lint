package golint

import (
	"github.com/surullabs/statictest"
	"os/exec"
	"strings"
)

// Check implements a golint Checker
type Check struct {
}

func (Check) Check(dir string) error {
	if err := statictest.InstallMissing("golint", "github.com/golang/lint/golint"); err != nil {
		return err
	}
	result, err := statictest.Exec(exec.Command("golint", dir))
	if err != nil {
		return err
	}
	str := strings.TrimSpace(result.Stdout)
	if str == "" {
		return nil
	}
	return &statictest.Error{Errors: strings.Split(str, "\n")}
}
