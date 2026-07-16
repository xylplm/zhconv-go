.PHONY: test vet demo build release-dry example

test:
	go test ./... -count=1

vet:
	go vet ./...

demo:
	go run ./cmd/zhconv -demo

example:
	cd examples/basic && go run .

build:
	go build -o bin/zhconv ./cmd/zhconv

release-dry:
	./scripts/build-release.sh v0.0.0-dev
