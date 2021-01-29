FROM alpine:3.12

ADD fio-exporter /bin/fio-exporter

RUN addgroup -g 777 exporter && adduser -u 777 -S -G exporter exporter \
    && apk add --no-cache fio

ENTRYPOINT ["/bin/fio-exporter"]

USER exporter:exporter
