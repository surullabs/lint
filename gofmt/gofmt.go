package gofmt

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/surullabs/statictest"
)

type Check struct {
}

func (Check) Check(pkg string) error {
	dir, err := statictest.PackageDir(pkg)
	if err != nil {
		return fmt.Errorf("gofmt-check: failed to read package dir: %v", err)
	}
	data, err := exec.Command("gofmt", "-d", dir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(data))
	}
	str := bytes.TrimSpace(data)
	if len(str) > 0 {
		return fmt.Errorf("File not formatted: %s", string(str))
	}
	return nil
}
