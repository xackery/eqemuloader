VERSION=$(shell git describe --tags)
LDFLAGS=-ldflags "-X main.Version=${VERSION}"
.PHONY: test
test:
	@(go test ./...)
vet:
	@(go vet ./...)
build:
	@(GOOS=windows go build -o dist/loader-windows.exe *.go)
	@(GOOS=linux go build -o dist/loader-linux *.go)
	@(GOOS=darwin go build -o dist/loader-darwin *.go)