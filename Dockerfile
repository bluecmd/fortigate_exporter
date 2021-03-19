# Use Prometheus' Golang Builder to avoid depending on Docker Hub
# See https://github.com/bluecmd/fortigate_exporter/issues/75 for more details
FROM quay.io/prometheus/golang-builder:1.16.2-base as builder

WORKDIR /build

COPY . .
RUN go get -v -t -d ./...
RUN CGO_ENABLED=0 go build -o main .

FROM scratch
WORKDIR /opt/fortigate_exporter

COPY --from=builder /build/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt .
ENV SSL_CERT_DIR=/opt/fortigate_exporter

EXPOSE 9710
ENTRYPOINT ["./main"]
CMD ["-auth-file", "/config/fortigate-key.yaml"]
