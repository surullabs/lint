// Package checkers provides utilities for implementing lint checkers.
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
)

// Error returns an error containing errs.
//
// If errs is empty, nil is returned. If not the returned
// error will implement the following interface
//
//     type errors interface {
//     	Errors() []string
//     }
//
// and return errs unmodified.
func Error(errs ...string) error {
	if len(errs) == 0 {
		return nil
	}
	return errorList(errs)
}

type errorList []string

func (e errorList) Errors() []string { return []string(e) }
func (e errorList) Error() string    { return strings.Join(e, "\n") }

func packageDir(path string) (string, error) {
	pkg, err := build.Import(path, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("import failed: %s: %v", path, err)
	}
	return pkg.Dir, nil
}

// FindBin returns bin if it exists in the path. If not it checks
// go bin directories ($GOROOT/bin and $GOPATH/bin) and returns that if it exists.
// If neither exist it returns an error.
func FindBin(bin string) (string, error) {
	if _, err := exec.LookPath(bin); err == nil {
		return bin, nil
	}
	srcDirs := build.Default.SrcDirs()
	for _, src := range srcDirs {
		binFile := filepath.Join(filepath.Dir(src), "bin", bin)
		if _, err := exec.LookPath(binFile); err == nil {
			return binFile, nil
		}
	}
	return "", fmt.Errorf("failed to find binary: %v", bin)
}

// InstallMissing runs go get getPath and then go get importPath
// if bin cannot be found in the directories contained in the PATH environment variable.
// It returns the path to the installed binary on success.
func InstallMissing(bin, getPath, importPath string) (string, error) {
	if b, err := FindBin(bin); err == nil {
		return b, nil
	}
	if data, err := exec.Command("go", "get", getPath).CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to get %s: %v: %s", importPath, err, string(data))
	}

	if data, err := exec.Command("go", "install", importPath).CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to install %s: %v: %s", importPath, err, string(data))
	}
	b, err := FindBin(bin)
	if err != nil {
		return "", fmt.Errorf("failed to lookup %v after install: %v", bin, err)
	}
	return b, nil
}

// ExecErrors holds errors collected from executing a linter
type ExecErrors []string

// Add appends any output from stdout or stderr as an error. The exit code is ignored.
// Empty lines are also ignored
func (e *ExecErrors) Add(r ExecResult) {
	str := strings.TrimSpace(
		strings.TrimSpace(r.Stdout) + "\n" + strings.TrimSpace(r.Stderr))
	if str == "" {
		return
	}
	*e = append(*e, strings.Split(str, "\n")...)
}

// Lint runs the linter specified by bin for each package in pkgs.
// The linter is installed if necessary using
//   go get getPath
//   go install installPath
//
// If getPath is empty, installPath is used for go get.
func Lint(bin, getPath, installPath string, pkgs []string, args ...string) error {
	if getPath == "" {
		getPath = installPath
	}
	b, err := InstallMissing(bin, getPath, installPath)
	if err != nil {
		return err
	}
	errs := &ExecErrors{}
	for _, pkg := range pkgs {
		p, perr := Load(pkg)
		if perr != nil {
			return fmt.Errorf("failed to load pkg info: %s: %v", pkg, perr)
		}
		result, _ := Exec(exec.Command(b, append(args, p.Path)...))
		errs.Add(result)
	}
	return Error((*errs)...)
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
	// build.Package instance for this package
	Build *build.Package
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

// SkipUnderscoreDirs skips all directories with the _ prefix
func SkipUnderscoreDirs(path, name string) bool { return strings.HasPrefix(name, "_") }

// SkipTestdata skips all testdata directories
func SkipTestdata(path, name string) bool { return name == "testdata" }

// SkipVendor skips all vendor directories
func SkipVendor(path, name string) bool { return name == "vendor" }

// SkipDirs returns a function that will return true if any of fns returns true.
func SkipDirs(fns ...func(path, name string) bool) func(string, string) bool {
	return func(p, n string) bool {
		for _, fn := range fns {
			if fn(p, n) {
				return true
			}
		}
		return false
	}
}

// SkipDirFunc determines if a directory must be skipped when listing packages.
// It takes two arguments, the full path and name of the directory and returns true
// if the directory should be skipped.
var SkipDirFunc = SkipDirs(SkipUnderscoreDirs, SkipTestdata)

func (p *Package) readPackages() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to find cwd: %v", err)
	}
	if filepath.Base(p.Path) != "..." {
		var b *build.Package
		if b, err = build.Import(p.Path, wd, build.FindOnly); err != nil {
			return fmt.Errorf("import failed: %s: %v", p.Path, err)
		}
		p.Pkgs, p.Build = []string{b.ImportPath}, b
		return nil
	}

	d := filepath.Dir(p.Path)
	b, err := build.Import(d, wd, build.FindOnly)
	if err != nil {
		return fmt.Errorf("import failed %s: %v", d, err)
	}
	p.Build = b
	dir := b.Dir
	var paths []string
	err = filepath.Walk(dir, func(path string, stat os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !stat.IsDir() {
			return nil
		}
		if SkipDirFunc(path, stat.Name()) {
			return filepath.SkipDir
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
