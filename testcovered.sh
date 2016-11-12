#!/bin/bash

set -e
set -x

go list -f '{{if len .XTestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile -coverpkg {{.ImportPath}},github.com/surullabs/lint/checkers {{.ImportPath}}"
{{end}}' ./... | xargs -L 1 sh -c
gover
