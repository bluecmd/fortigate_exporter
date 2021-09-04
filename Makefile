VERSION := $(shell git describe --tags)
GIT_HASH := $(shell git rev-parse --short HEAD )

GO_VERSION        ?= $(shell go version)
GO_VERSION_NUMBER ?= $(word 3, $(GO_VERSION))
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.GitHash=${GIT_HASH} -X main.GoVersion=${GO_VERSION_NUMBER}"

.PHONY: build
build:
	CGO_ENABLED=0 go build ${LDFLAGS} -v -o target/fortigate-exporter .

.PHONY: build-release
build-release: build-release-amd64 build-release-arm64

.PHONY: build-release-amd64
build-release-amd64:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build ${LDFLAGS} -o=fortigate-exporter.linux.amd64 .        && \
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o=fortigate-exporter.windows.amd64.exe .  && \
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build ${LDFLAGS} -o=fortigate-exporter.darwin.amd64 .

.PHONY: build-release-arm64
build-release-arm64:
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build ${LDFLAGS} -o=fortigate-exporter.linux.arm64 .        && \
 	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build ${LDFLAGS} -o=fortigate-exporter.darwin.arm64 .

.PHONY: generate-html-coverage
generate-html-coverage:
	go tool cover -html=cover.out -o coverage.html
	@printf "Generated coverage html \n"

.PHONY: print-coverage
print-coverage:
	@go tool cover -func cover.out

.PHONY: test-unittests
test-unittests:
	go test -v -race -coverprofile cover.out ./...

.PHONY: test
test: fmt-check vet test-unittests generate-html-coverage print-coverage
	@printf "Sucessfully run tests \n"

.PHONY: get-dependencies
get-dependencies:
	go get -v -t -d ./...

.PHONY: vet
vet:
	@go vet ./...
	@go mod tidy

test-output:
	$(shell echo $$GO_VERSION_NUMBER)

.PHONY: fmt-fix
fmt-fix:
	@go mod download golang.org/x/tools
	@go run golang.org/x/tools/cmd/goimports -w -l .

.PHONY: fmt-check
fmt-check:
	@printf "Check formatting... \n"
	@go mod download golang.org/x/tools
	@if [[ $$( go run golang.org/x/tools/cmd/goimports -l . ) ]]; then printf "Files not properly formatted. Run 'make fmt-fix' \n"; exit 1; else printf "Check formatting finished \n"; fi
