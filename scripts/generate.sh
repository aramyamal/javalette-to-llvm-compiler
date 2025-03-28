#!/usr/bin/env bash

set -e

echo "Running code generation script:"

SCRIPT_DIR=$(dirname "$(realpath $0)")
ANTLR_JAR="antlr-4.13.2-complete.jar"
ANTLR_JAR_DIR="$SCRIPT_DIR/../tools/"
ANTLR_JAR_PATH="$ANTLR_JAR_DIR$ANTLR_JAR"
GEN_DIR="$SCRIPT_DIR/../gen/parsing/"

mkdir -p "$ANTLR_JAR_DIR"
mkdir -p "$GEN_DIR"

if [[ ! -f "$ANTLR_JAR_PATH" ]]; then
    echo "Downloading antlr4 jar into tools/..."
    curl -O "https://www.antlr.org/download/$ANTLR_JAR" --output-dir "$ANTLR_JAR_DIR"
else
    echo "antlr4 jar already in tools/, reusing it"
fi

echo "Generating parser code with antlr4 from grammar file to gen/..."
java -Xmx500M -cp "$ANTLR_JAR_PATH:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Go -no-visitor -package parsing *.g4 -o "$GEN_DIR"

echo "Finished generating parser code"
