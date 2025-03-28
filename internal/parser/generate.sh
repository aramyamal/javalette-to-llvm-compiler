#!/usr/bin/env bash

echo "Running code generation script:"

ANTLR_JAR="antlr-4.13.2-complete.jar"
ANTLR_JAR_DIR="../../tools/antlr4/"
ANTLR_JAR_PATH="$ANTLR_JAR_DIR$ANTLR_JAR"
GEN_DIR="../../gen/parsing/"

mkdir -p "$ANTLR_JAR_DIR"
mkdir -p "$GEN_DIR"

if [[ ! -f "$ANTLR_JAR_PATH" ]]; then
    echo "Downloading antlr4 jar..."
    curl -O "https://www.antlr.org/download/$ANTLR_JAR" --output-dir "$ANTLR_JAR_DIR"
else
    echo "antlr4 jar already downloaded before, reusing it"
fi

echo "Generating parser code with antlr4 from grammar file..."
java -Xmx500M -cp "$ANTLR_JAR_PATH:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Go -no-visitor -package parsing *.g4 -o "$GEN_DIR"

