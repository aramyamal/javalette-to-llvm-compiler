SOURCES  := $(shell find . -name "*.go")

all: clean generate build

generate:
	@echo "Generating parser..."
	@go generate ./...

build: $(SOURCES)
	@bash scripts/build.sh

clean:
	@bash scripts/clean.sh

.PHONY: all clean build generate
