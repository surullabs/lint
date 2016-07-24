package gofmt

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Check struct {
}

func (Check) Check(pkg string) error {
	data, err := exec.Command("go", "fmt", pkg).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(data))
	}
	str := bytes.TrimSpace(data)
	if len(str) > 0 {
		return fmt.Errorf("File not formatted: %s", string(str))
	}
	return nil
}
