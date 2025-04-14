#!/usr/bin/env bash
 
set -e

echo "Building jlc executable..."

SCRIPT_DIR=$(dirname "$(realpath $0)")

cd "$SCRIPT_DIR/.."

mkdir -p build

go build -o build/jlc cmd/jlc/main.go 

echo "Finished building"

