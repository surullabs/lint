package gometalinter

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"os/exec"

	"github.com/surullabs/lint/checkers"
)

// Check implements a check using a vendored version of gometalinter. Args are the
// arguments passed to gometalinter. Do not include directory names in Args. These
// will be added automatically, based on the arguments to Check(pkgs).
type Check struct {
	Args []string
}

// Check runs a vendored version of gometalinter. It builds the
// metalinter by detecting the location of the vendor directory and
// using that as the GOPATH for building the metalinter binary. This is
// similar to what gometalinter does internally.
func (c Check) Check(pkgs ...string) error {
	dirs := make([]string, len(pkgs))
	for i, pkg := range pkgs {
		p, err := checkers.Load(pkg)
		if err != nil {
			return err
		}
		dirs[i] = p.Build.Dir
		if filepath.Base(pkg) == "..." {
			dirs[i] = filepath.Join(dirs[i], "...")
		}
	}
	return runMetalinter(append(c.Args, dirs...)...)
}

func runMetalinter(args ...string) error {
	env, bin, err := installMetaLinter()
	if err != nil {
		return err
	}
	cmd := exec.Command(bin, args...)
	cmd.Env = env
	r, err := checkers.Exec(cmd)
	// From the gometalinter README it sets two bits of information in the error code.
	// So any error code from 1 - 3 is a metalinter error which we pass on. Any other
	// code is an unexpected error so return a regular error.
	if err != nil && (r.Code < 0 || r.Code > 3) {
		return fmt.Errorf("gometalinter exec failed: %v:\n%s\n%s", err, r.Stdout, r.Stderr)
	}
	res := &checkers.ExecErrors{}
	res.Add(r)
	return checkers.Error((*res)...)
}

func installMetaLinter() ([]string, string, error) {
	root := ""
	// Look up the actual package path of the lint install instead of assuming it is
	// github.com/surullabs/lint. This is needed to handle cases where this library
	// is itself vendored.
	path := filepath.Dir(reflect.TypeOf(Check{}).PkgPath())
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return nil, "", fmt.Errorf("%s: GOPATH not set when looking for source location", path)
	}
	gopaths := strings.Split(gopath, string(os.PathListSeparator))
	for _, p := range gopaths {
		lintPath := filepath.Join(p, "src", path)
		if _, err := os.Stat(lintPath); err == nil {
			root = filepath.Join(lintPath, "gometalinter", "_vendored")
			break
		}
	}
	if root == "" {
		return nil, "", fmt.Errorf("%s: not found under GOPATH (%v)", path, gopath)
	}

	// Check to see if the bin exists
	bin := filepath.Join(root, "bin", "gometalinter")
	// Always rebuild the binary. This might seem wasteful, but if the vendored version
	// is updated we need to have the latest version of the binary rebuilt.
	env := os.Environ()
	gi, orig := -1, ""
	for i := range env {
		// Messing with the local process environment can have undesirable effects, so
		// create a copy and replace GOPATH
		if strings.HasPrefix(env[i], "GOPATH=") {
			env[i], gi, orig = fmt.Sprintf("GOPATH=%v", root), i, env[i]
		} else if strings.HasPrefix(env[i], "PATH=") {
			env[i] = fmt.Sprintf("PATH=%s%c%s", filepath.Join(root, "bin"), filepath.ListSeparator, env[i])
		}
	}
	cmd := exec.Command("go", "install", "github.com/alecthomas/gometalinter")
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, "", fmt.Errorf("failed to install gometalinter: %v\n%s", err, string(out))
	}
	if _, err := os.Stat(bin); err != nil {
		return nil, "", fmt.Errorf("gometalinter not installed at %v: %v", bin, err)
	}
	cmd = exec.Command(bin, "--install")
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, "", fmt.Errorf("failed to install vendored linters: %v\n%s", err, string(out))
	}
	env[gi] = orig
	return env, bin, nil
}
