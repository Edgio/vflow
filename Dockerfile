# build vFlow in the first stage
FROM golang:1.8.3 as builder

WORKDIR /go/src/

RUN mkdir -p ./github.com/VerizonDigital/vflow
ADD . ./github.com/VerizonDigital/vflow
WORKDIR ./github.com/VerizonDigital/vflow/vflow

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vflow .

# run vFlow within alpine in the second stage
FROM alpine:latest
COPY --from=builder /go/src/github.com/VerizonDigital/vflow/vflow/vflow /usr/bin/vflow
COPY ./scripts/dockerStart.sh /dockerStart.sh

EXPOSE 4739 6343 9996 4729 8081

VOLUME /etc/vflow

CMD sh /dockerStart.sh
