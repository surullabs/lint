package dupl

import (
	"bytes"
	"fmt"
	"os/exec"

	"strconv"

	"regexp"
	"strings"

	"github.com/sridharv/ternary"
	"github.com/surullabs/statictest"
	"github.com/surullabs/statictest/checkers"
)

// Check is implements statictest.Checker for gofmt.
type Check struct {
	Threshold int
}

var (
	foundRE = regexp.MustCompile(`found [0-9]+ clones:`)
	finalRE = regexp.MustCompile(`Found total [0-9]+ clone groups.`)
)

// SkipTwo implements Skipper and skips all duplicates having just two instances.
var SkipTwo = statictest.SkipFunc(func(str string) bool {
	return strings.HasPrefix(str, "dupl.Check: found 2 clones:") ||
		strings.HasPrefix(str, "found 2 clones:")
})

// Skip returns a Skipper which ignores the dup whose first instance ends with suffix.
//
//    c = statictest.Skip(c, Skip("lint.go:1,12"))
func Skip(suffix string) statictest.Skipper {
	return statictest.SkipFunc(func(str string) bool {
		if !strings.Contains(str, "dupl") {
			return false
		}
		lines := strings.Split(str, "\n")
		return len(lines) > 1 && strings.HasSuffix(lines[1], suffix)
	})
}

// Check runs
//   dupl <files>
//
// for all files in pkgs.
func (c Check) Check(pkgs ...string) error {
	files, err := checkers.GoFiles(pkgs...)
	if err != nil {
		return err
	}
	if err := checkers.InstallMissing("dupl", "github.com/mibk/dupl"); err != nil {
		return err
	}
	t := ternary.Int(c.Threshold == 0, 15, c.Threshold)
	args := append([]string{"-t", strconv.Itoa(t)}, files...)
	data, err := exec.Command("dupl", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("dupl failed: %v: %s", err, string(data))
	}
	data = bytes.TrimSpace(data)
	loc := finalRE.FindIndex(data)
	if loc == nil {
		return fmt.Errorf("unexpected output: couldn't find final clone group line: %v", string(data))
	}
	if loc[0] == 0 {
		return nil
	}
	data = data[0:loc[0]]
	indices := foundRE.FindAllIndex(data, -1)
	if indices == nil {
		return fmt.Errorf("%s", string(data))
	}
	indices = append(indices, []int{len(data), len(data)})
	p := -1
	errs := &statictest.Error{}
	for _, l := range indices {
		if p >= 0 {
			errs.Errors = append(errs.Errors, string(data[p:l[0]]))
		}
		p = l[0]
	}
	return errs.AsError()
}
