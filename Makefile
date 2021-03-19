VERSION := $(shell git tag | grep ^v | sort -V | tail -n 1)
GIT_HASH := $(shell git rev-parse --short HEAD )
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.GitHash=${GIT_HASH}"

build:
	go build ${LDFLAGS} -v -o target/fortigate-exporter .

build-release:
	GOOS=linux   GOARCH=amd64 go build ${LDFLAGS} -o=target/fortigate-exporter.linux.amd64 .        && \
  	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o=target/fortigate-exporter.windows.amd64.exe .  && \
  	GOOS=darwin  GOARCH=amd64 go build ${LDFLAGS} -o=target/fortigate-exporter.darwin.amd64 .

test:
	go test -v .

get-dependencies:
	go get -v -t -d ./...

vet:
	go vet ./...
