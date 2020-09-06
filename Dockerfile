FROM golang:latest

WORKDIR /opt/fortigate_exporter

COPY . .
RUN go get -v -t -d ./...

RUN go build -o main .

EXPOSE 9710
CMD ["./main", "-auth-file", "/opt/fortigate-key.yaml"]