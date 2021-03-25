VERSION := $(shell git describe --tags)
GIT_HASH := $(shell git rev-parse --short HEAD )

GO_VERSION        ?= $(shell go version)
GO_VERSION_NUMBER ?= $(word 3, $(GO_VERSION))
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.GitHash=${GIT_HASH} -X main.GoVersion=${GO_VERSION_NUMBER}"

.PHONY: build
build:
	go build ${LDFLAGS} -v -o target/fortigate-exporter .

.PHONY: build-release
build-release: build-release-amd64 build-release-arm64

.PHONY: build-release-amd64
build-release-amd64:
	GOOS=linux   GOARCH=amd64 go build ${LDFLAGS} -o=fortigate-exporter.linux.amd64 .        && \
  	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o=fortigate-exporter.windows.amd64.exe .  && \
  	GOOS=darwin  GOARCH=amd64 go build ${LDFLAGS} -o=fortigate-exporter.darwin.amd64 .

.PHONY: build-release-arm64
build-release-arm64:
	GOOS=linux   GOARCH=arm64 go build ${LDFLAGS} -o=fortigate-exporter.linux.arm64 .        && \
  	GOOS=darwin  GOARCH=arm64 go build ${LDFLAGS} -o=fortigate-exporter.darwin.arm64 .

.PHONY: test
test:
	go test -v .

.PHONY: get-dependencies
get-dependencies:
	go get -v -t -d ./...

.PHONY: vet
vet:
	go vet ./...

test-output:
	$(shell echo $$GO_VERSION_NUMBER)