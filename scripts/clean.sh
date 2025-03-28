#!/usr/bin/env bash

set -e

echo "Cleaning up all build artefacts..."

SCRIPT_DIR=$(dirname "$(realpath $0)")

cd "$SCRIPT_DIR/.."

echo "Removing bin/"
rm -rf ./bin 
echo "Removing gen/"
rm -rf ./gen
echo "Removing antlr4 jar from tools/"
rm -rf ./tools/

echo "Cleanup completed"
