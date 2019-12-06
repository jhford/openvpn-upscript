#!/bin/bash
set -e

export GOOS=linux
export GOARCH=arm
export GOARM=7

source build.sh
