.PHONY: all clean build generate

all: generate build

generate:
	@go generate ./...

build:
	@bash scripts/build.sh

clean:
	@bash scripts/clean.sh

