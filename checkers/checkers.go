// Package checkers provides utilities for implementing statictest checkers.
package checkers

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/surullabs/statictest"
)

// PackageDir returns the directory containing a package.
func PackageDir(path string) (string, error) {
	pkg, err := build.Import(path, ".", build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.Dir, nil
}

// InstallMissing runs go get importPath if bin cannot be found in the directories
// contained in the PATH environment variable.
func InstallMissing(bin, importPath string) error {
	if _, err := exec.LookPath(bin); err == nil {
		return nil
	}
	if data, err := exec.Command("go", "get", importPath).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install %s: %v: %s", importPath, err, string(data))
	}
	return nil
}

// Lint runs the linter specified by bin for pkg, installing it if necessary using
// go get importPath.
func Lint(bin, importPath, pkg string) error {
	if err := InstallMissing(bin, importPath); err != nil {
		return err
	}
	result, _ := Exec(exec.Command(bin, pkg))
	str := strings.TrimSpace(
		strings.TrimSpace(result.Stdout) + "\n" + strings.TrimSpace(result.Stderr))
	if str == "" {
		return nil
	}
	return &statictest.Error{Errors: strings.Split(str, "\n")}
}

// ExecResult holds a status code, stdout and stderr for a single command execution.
type ExecResult struct {
	Code   int
	Stdout string
	Stderr string
}

// Exec executes cmd and results exit code, stdout and stderr of the result. It
// additionally returns and error if the status code is not 0.
func Exec(cmd *exec.Cmd) (ExecResult, error) {
	res := ExecResult{Code: -1}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return res, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return res, err
	}
	if err = cmd.Start(); err != nil {
		return res, err
	}
	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		return res, fmt.Errorf("failed to read stdout: %v", err)
	}
	res.Stdout = string(data)
	if data, err = ioutil.ReadAll(stderr); err != nil {
		return res, fmt.Errorf("failed to read stderr: %v", err)
	}
	res.Stderr = string(data)
	err = cmd.Wait()
	if st, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		res.Code = st.ExitStatus()
	}
	return res, err
}

// Files returns all files in a package
func Files(pkg string) ([]string, error) {
	dir, err := PackageDir(pkg)
	if err != nil {
		return nil, err
	}
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files, i := make([]string, len(entries)), 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files[i] = filepath.Join(dir, entry.Name())
		i++
	}
	return files[:i], nil
}

// GoFiles returns all .go files in a package.
func GoFiles(pkg string) ([]string, error) {
	files, err := Files(pkg)
	if err != nil {
		return nil, err
	}
	gofiles, i := make([]string, len(files)), 0
	for _, f := range files {
		if !strings.HasSuffix(f, ".go") {
			continue
		}
		gofiles[i] = f
		i++
	}
	return gofiles[:i], nil
}
