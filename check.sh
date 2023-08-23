#!/bin/sh -eu

buf lint
buf format -d
buf breaking --against '.git#branch=main,recurse_submodules=true'

yamlfmt -lint
yamllint .

golangci-lint run
