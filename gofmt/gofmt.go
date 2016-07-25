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
	files, err := statictest.GoFiles(pkg)
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
