## Lint - run linters from Go

[![Build Status](https://travis-ci.org/surullabs/lint.svg?branch=master)](https://travis-ci.org/surullabs/lint) [![GoDoc](https://godoc.org/github.com/surullabs/lint?status.svg)](https://godoc.org/github.com/surullabs/lint) [![Coverage Status](https://coveralls.io/repos/github/surullabs/lint/badge.svg?branch=master)](https://coveralls.io/github/surullabs/lint?branch=master)

Lint makes it easy to run linters from Go code. This allows lint checks to be part of a regular `go build` + `go test` workflow. False positives are easily ignored and linters are automatically integrated into CI pipelines without any extra effort. Check the [project website](https://www.timeferret.com/lint) to learn more about how it can be useful.

### Quick Start

Download using
```
go get -t github.com/surullabs/lint
```
Run the default linters by adding a new test at the top level of your repository
```
func TestLint(t *testing.T) {
    // Run default linters
    err := lint.Default.Check("./...")
    
    // Ignore lint errors from auto-generated files
    err = lint.Skip(err, lint.RegexpMatch(`_string\.go`, `\.pb\.go`))
    
    if err != nil {
        t.Fatal("lint failures: %v", err)
    }
}
```

### How it works

`lint` runs linters using the excellent `os/exec` package. It searches all Go binary directories for the needed binaries and when they don't exist it downloads them using `go get`. Errors generated by running linters are split by newline and can be skipped as needed.

### Default linters

  - `gofmt` - [Run `gofmt -d` and report any differences as errors](https://golang.org/cmd/gofmt/)
  - `govet` - [Run `go tool vet -shadow`](https://golang.org/cmd/vet/)
  - `golint` - [https://github.com/golang/lint](https://github.com/golang/lint)
  - `gosimple` - [Code simplification](https://github.com/dominikh/go-simple)
  - `gostaticcheck` - [Verify function arguments](https://github.com/dominikh/go-staticcheck)
  - `errcheck` - [Find ignored errors](https://github.com/kisielk/errcheck)
  
 ## Using `gometalinter`

[Gometalinter](https://github.com/alecthomas/gometalinter) runs a number of linters concurrently. It also vendors each of these and uses the vendored versions automatically. A vendored version of `gometalinter` is included and can be used in the following manner. Please note that not all linters used by gometalinter have been tested.
 
```
import (
    "testing"
    "github.com/surullabs/lint/gometalinter"
)

func TestLint(t *testing.T) {
    // Run default linters
    metalinter := gometalinter.Check{
        Args: []string{
            // Arguments to gometalinter. Do not include the package names here.
        },
    }
    if err := metalinter.Check("./..."); err != nil {
        t.Fatal("lint failures: %v", err)
    }
}

```
 
 ## Other available linters
 
  - `varcheck` - [Detect unused variables and constants](https://github.com/opennota/check)
  - `structcheck` - [Detect unused struct fields](https://github.com/opennota/check)
  - `aligncheck` - [Detect suboptimal struct alignment](https://github.com/opennota/check)
  - `dupl` - [Detect duplicated code](https://github.com/mibk/dupl)
 
### Why `lint`?

There are a number of excellent linters available for Go and Lint makes it easy to run them from tests. While building our mobile calendar app [TimeFerret](https://www.timeferret.com), (which is built primarily in Go), including scripts that run linters as part of every repository grew tiresome very soon. Using `lint` to create tests that ran on each commit made the codebase much more stable, since any unneeded false positives were easily skipped. The main advantages of using `lint` over running tools manually is:

  - Skip false positives explicitly in your tests - This makes it easy to run only needed checks.
  - Enforce linter usage with no overhead - No special build scripts are needed to install linters on each developer machine as they are automatically downloaded.
  - Simple CI integration - Since linters are run as part of tests, there are no extra steps needed to integrate them into your CI pipeline.

### Adding a custom linter

Adding a new linter is made dead simple by the `github.com/surullabs/lint/checkers` package. The entire source for the `golint` integration is

```
import "github.com/surullabs/lint/checkers"

type Check struct {
}

func (Check) Check(pkgs ...string) error {
    return checkers.Lint("golint", "", github.com/golang/lint/golint", pkgs)
}
```

The `github.com/surullabs/lint/testutil` package contains utilities for testing custom linters.

You can also take a look at [this CL](https://github.com/surullabs/lint/commit/5e6be15e3b9964e8465655abb9759defd1c46af9) which adds `varcheck` for an example of how to add a linter.
### License

Lint is available under the Apache License. See the LICENSE file for details.

### Contributing

Pull requests are always welcome! Please ensure any changes you send have an accompanying test case.