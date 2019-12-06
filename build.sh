#!/bin/bash
set -e

version="0.0.1"
commit="$(git rev-parse HEAD)"

go build -ldflags "-X pkg.Version=$version -X pkg.Commit=$commit" ./cmd/up