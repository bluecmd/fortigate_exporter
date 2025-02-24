ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:glibc
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/fortigate_exporter /bin/fortigate_exporter

EXPOSE      9710
USER        nobody
ENTRYPOINT  [ "/bin/fortigate_exporter" ]
