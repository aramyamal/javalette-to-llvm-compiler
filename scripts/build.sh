#!/usr/bin/env bash
 
set -e


SCRIPT_DIR=$(dirname "$(realpath $0)")

cd "$SCRIPT_DIR/.."

mkdir -p build

echo "Building jlc executable..."
go build -o build/jlc cmd/jlc/main.go
echo "Building typecheck executable..."
go build -o build/typecheck cmd/typecheck/main.go

echo "Finished building"

