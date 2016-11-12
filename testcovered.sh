#!/bin/bash

set -e
set -x

go list -f '{{if len .XTestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"
{{end}}' ./... | xargs -L 1 sh -c
rm ./.coverprofile
go test -v -coverprofile=./.coverprofile -coverpkg github.com/surullabs/lint,github.com/surullabs/lint/checkers
gover
