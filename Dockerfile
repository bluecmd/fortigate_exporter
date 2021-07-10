# Use Prometheus' Golang Builder to avoid depending on Docker Hub
# See https://github.com/bluecmd/fortigate_exporter/issues/75 for more details
FROM quay.io/prometheus/golang-builder:1.16.2-base as builder

WORKDIR /build

COPY . .
RUN go get -v -t -d ./...
RUN make build

FROM scratch
WORKDIR /opt/fortigate_exporter

COPY --from=builder /build/target/fortigate-exporter .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt .
ENV SSL_CERT_DIR=/opt/fortigate_exporter

EXPOSE 9710
ENTRYPOINT ["./fortigate-exporter"]
CMD ["-auth-file", "/config/fortigate-key.yaml"]
