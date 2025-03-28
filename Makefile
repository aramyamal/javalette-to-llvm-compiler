.PHONY: all clean build generate

generate:
	@go generate ./...

build:
	@bash scripts/build.sh

clean:
	@bash scripts/clean.sh

all: generate build
