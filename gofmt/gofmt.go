package gofmt

import (
	"bytes"
	"fmt"
	"github.com/surullabs/lint/checkers"
	"os/exec"
)

// Check is implements lint.Checker for gofmt.
type Check struct {
}

// Check runs
//   gofmt -d <files>
//
// for all files in pkgs.
func (Check) Check(pkgs ...string) error {
	files, err := checkers.GoFiles(pkgs...)
	if err != nil {
		return err
	}
	data, err := exec.Command("gofmt", append([]string{"-d"}, files...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(data))
	}
	str := bytes.TrimSpace(data)
	if len(str) > 0 {
		return fmt.Errorf("File not formatted: %s", string(str))
	}
	return nil
}
