package dupl

import (
	"bytes"
	"fmt"
	"os/exec"

	"strconv"

	"regexp"
	"strings"

	"github.com/surullabs/lint"
	"github.com/surullabs/lint/checkers"
)

// Check is implements lint.Checker for gofmt.
type Check struct {
	Threshold int
}

var (
	foundRE = regexp.MustCompile(`found [0-9]+ clones:`)
	finalRE = regexp.MustCompile(`Found total [0-9]+ clone groups.`)
)

type skipFunc func(str string) bool

func (s skipFunc) Skip(str string) bool { return s(str) }

// SkipTwo implements Skipper and skips all duplicates having just two instances.
var SkipTwo = skipFunc(func(str string) bool {
	return strings.HasPrefix(str, "dupl.Check: found 2 clones:") ||
		strings.HasPrefix(str, "found 2 clones:")
})

// Skip returns a Skipper which ignores the dup whose first instance ends with suffix.
//
//    c = lint.Skip(c, Skip("lint.go:1,12"))
func Skip(suffix string) lint.Skipper {
	return skipFunc(func(str string) bool {
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
	bin, err := checkers.InstallMissing("dupl", "github.com/mibk/dupl", "github.com/mibk/dupl")
	if err != nil {
		return err
	}
	t := c.Threshold
	if t == 0 {
		t = 15
	}
	args := append([]string{"-t", strconv.Itoa(t)}, files...)
	data, err := exec.Command(bin, args...).CombinedOutput()
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
	var errs []string
	for _, l := range indices {
		if p >= 0 {
			errs = append(errs, string(data[p:l[0]]))
		}
		p = l[0]
	}
	return checkers.Error(errs...)
}
