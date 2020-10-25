FROM golang:latest as builder

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
