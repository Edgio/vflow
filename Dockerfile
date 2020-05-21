# build vFlow in the first stage
FROM golang:1.14.3 as builder
ENV VERSION 0.7.0
WORKDIR /go/src/

RUN mkdir -p github.com/VerizonDigital/vflow
ADD . github.com/VerizonDigital/vflow
WORKDIR github.com/VerizonDigital/vflow/vflow
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=$VERSION" -a -installsuffix cgo .

# run vFlow within alpine in the second stage
FROM alpine:latest
COPY --from=builder /go/src/github.com/VerizonDigital/vflow/vflow/vflow /usr/bin/
COPY scripts/dockerStart.sh /dockerStart.sh

EXPOSE 4739 6343 9996 4729 8081

VOLUME /etc/vflow

CMD sh /dockerStart.sh
