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

	"sync"

	"os"

	"github.com/surullabs/statictest"
)

func packageDir(path string) (string, error) {
	pkg, err := build.Import(path, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("import failed: %s: %v", path, err)
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

// Lint runs the linter specified by bin for each package in pkgs.
// The linter is installed if necessary using go get importPath.
func Lint(bin, importPath string, pkgs []string) error {
	if err := InstallMissing(bin, importPath); err != nil {
		return err
	}
	var errs []string
	for _, pkg := range pkgs {
		p, err := Load(pkg)
		if err != nil {
			return fmt.Errorf("failed to load pkg info: %s: %v", pkg, err)
		}
		result, _ := Exec(exec.Command(bin, p.Path))
		str := strings.TrimSpace(
			strings.TrimSpace(result.Stdout) + "\n" + strings.TrimSpace(result.Stderr))
		if str == "" {
			continue
		}
		errs = append(errs, strings.Split(str, "\n")...)
	}
	if len(errs) == 0 {
		return nil
	}
	return &statictest.Error{Errors: errs}
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
		return res, fmt.Errorf("cmd.StdoutPipe failed: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return res, fmt.Errorf("cmd.StderrPipe failed: %v", err)
	}
	if err = cmd.Start(); err != nil {
		return res, fmt.Errorf("cmd.Start failed: %v", err)
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

// GoFiles lists all .go files in pkgs.
func GoFiles(pkgs ...string) ([]string, error) {
	var files []string
	for _, pkg := range pkgs {
		p, err := Load(pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to load go files for %s: %v", pkg, err)
		}
		files = append(files, p.GoFiles...)
	}
	return files, nil
}

func filterGoFiles(files []string) []string {
	gofiles, i := make([]string, len(files)), 0
	for _, f := range files {
		if !strings.HasSuffix(f, ".go") || strings.HasSuffix(f, "_test.go") {
			continue
		}
		gofiles[i] = f
		i++
	}
	return gofiles[:i]
}

// Package holds file and sub directory information about a single package.
type Package struct {
	// The path of the package. Can be a wildcard path such as ...
	Path string
	// Files holds all files in the package. If the Path is a wildcard path
	// (...) files in sub packages are also returned.
	Files []string
	// All files in Files with a .go extension
	GoFiles []string
	// All sub packages if Path is a wildcard, or just Path if not.
	Pkgs []string
}

var (
	packages     = map[string]*Package{}
	packageMutex sync.Mutex
)

// Unload removes any information about pkg from the cache.
func Unload(pkg string) {
	packageMutex.Lock()
	defer packageMutex.Unlock()
	delete(packages, pkg)
}

// Load returns a cached Package instance if one exists or creates a new instance if not.
// It returns an error if there was an error reading package information.
func Load(pkg string) (*Package, error) {
	packageMutex.Lock()
	defer packageMutex.Unlock()
	p := packages[pkg]
	if p != nil {
		return p, nil
	}
	p = &Package{Path: pkg}
	if err := p.load(); err != nil {
		return nil, err
	}
	packages[pkg] = p
	return p, nil
}

func (p *Package) load() error {
	if err := p.readPackages(); err != nil {
		return fmt.Errorf("failed read sub packages: %s: %v", p.Path, err)
	}
	if err := p.readFiles(); err != nil {
		return fmt.Errorf("failed to list files: %s: %v", p.Path, err)
	}
	p.GoFiles = filterGoFiles(p.Files)
	return nil
}

func (p *Package) readFiles() error {
	var res []string
	for _, pkg := range p.Pkgs {
		dir, err := packageDir(pkg)
		if err != nil {
			return err
		}
		entries, err := ioutil.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to list dir %s: %v", dir, err)
		}
		files, i := make([]string, len(entries)), 0
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			files[i] = filepath.Join(dir, entry.Name())
			i++
		}
		res = append(res, files[:i]...)
	}
	p.Files = res
	return nil
}

func (p *Package) readPackages() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to find cwd: %v", err)
	}
	if filepath.Base(p.Path) != "..." {
		b, err := build.Import(p.Path, wd, build.FindOnly)
		if err != nil {
			return fmt.Errorf("import failed: %s: %v", p.Path, err)
		}
		p.Pkgs = []string{b.ImportPath}
		return nil
	}

	d := filepath.Dir(p.Path)
	b, err := build.Import(d, wd, build.FindOnly)
	if err != nil {
		return fmt.Errorf("import failed %s: %v", d, err)
	}
	dir := b.Dir
	var paths []string
	err = filepath.Walk(dir, func(path string, stat os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !stat.IsDir() {
			return nil
		}
		p, perr := build.ImportDir(path, build.FindOnly)
		if perr != nil {
			if _, noGo := perr.(*build.NoGoError); noGo {
				return nil
			}
			return fmt.Errorf("import failed: %s: %v", path, perr)
		}
		paths = append(paths, p.ImportPath)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to list %s: %v", dir, err)
	}
	p.Pkgs = paths
	return nil
}
