#!/usr/bin/env bash

set -e
shopt -s extglob

echo "Cleaning up all build artefacts..."

SCRIPT_DIR=$(dirname "$(realpath $0)")

cd "$SCRIPT_DIR/.."

echo "Removing build/"
rm -rf ./build 
echo "Removing gen/"
cd ./gen/parser/
rm -rf !(doc.go)
# echo "Removing antlr4 jar from tools/"
# rm -rf ./tools/

echo "Cleanup completed"
