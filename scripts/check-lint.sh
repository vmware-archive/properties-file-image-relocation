#!/usr/bin/env bash
golangci-lint -E misspell,gocyclo,dupl,gofmt,golint,unconvert,goimports,depguard,gocritic,interfacer run --disable-all
