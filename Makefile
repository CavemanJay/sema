VERSION=$(shell git describe --tags --abbrev=0)

VERSION_BUILD_FLAG := -ldflags '-X main.version=$(VERSION)'

build:
	go build -o ./bin/ $(VERSION_BUILD_FLAG) ./cmd/sema

run: build
	./bin/sema

install:
# go install -ldflags -w -s -X main.version=$(VERSION)
	go install -ldflags '-w -s' $(VERSION_BUILD_FLAG) ./cmd/sema