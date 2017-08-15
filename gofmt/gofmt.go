package gofmt

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/surullabs/lint/checkers"
)

// Check is implements lint.Checker for gofmt.
type Check struct {
}

// Check runs
//   gofmt -d <files>
//
// for all files in pkgs.
func (Check) Check(pkgs ...string) error {
	var errs = []string{}

	files, err := checkers.GoFiles(pkgs...)
	if err != nil {
		return err
	}

	for _, f := range files {
		data, err := exec.Command("gofmt", "-d", f).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %s", err, string(data))
		}

		str := bytes.TrimSpace(data)
		if len(str) > 0 {
			errs = append(errs, fmt.Sprintf("File not formatted: %s", string(str)))
		}
	}

	return checkers.Error(errs...)
}
