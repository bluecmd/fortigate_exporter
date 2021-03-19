VERSION := $(shell git describe --tags)
GIT_HASH := $(shell git rev-parse --short HEAD )

GO_VERSION        ?= $(shell go version)
GO_VERSION_NUMBER ?= $(word 3, $(GO_VERSION))
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.GitHash=${GIT_HASH} -X main.GoVersion=${GO_VERSION_NUMBER}"

.PHONY: build
build:
	go build ${LDFLAGS} -v -o target/fortigate-exporter .

.PHONY: build-release
build-release:
	GOOS=linux   GOARCH=amd64 go build ${LDFLAGS} -o=target/fortigate-exporter.linux.amd64 .        && \
  	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o=target/fortigate-exporter.windows.amd64.exe .  && \
  	GOOS=darwin  GOARCH=amd64 go build ${LDFLAGS} -o=target/fortigate-exporter.darwin.amd64 .

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