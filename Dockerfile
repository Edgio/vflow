# build vFlow in the first stage
FROM golang:1.15.3 as builder
WORKDIR /go/src/

RUN mkdir -p github.com/EdgeCast/vflow
ADD . github.com/EdgeCast/vflow
WORKDIR /go/src/github.com/EdgeCast/vflow
RUN make build

# run vFlow within alpine in the second stage
FROM alpine:latest
COPY --from=builder /go/src/github.com/EdgeCast/vflow/vflow/vflow /usr/bin/
COPY scripts/dockerStart.sh /dockerStart.sh

EXPOSE 4739 6343 9996 4729 8081

VOLUME /etc/vflow

CMD sh /dockerStart.sh
