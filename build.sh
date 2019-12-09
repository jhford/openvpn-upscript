#!/bin/bash
set -e

version="0.0.1"
commit="$(git rev-parse HEAD)"

echo $version $commit

go build -ldflags "-s -w -X github.com/jhford/openvpn-up/pkg.Version=$version -X github.com/jhford/openvpn-up/pkg.Commit=$commit" ./cmd/up
upx --best --brute ./up
