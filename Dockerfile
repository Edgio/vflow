FROM ubuntu:trusty

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    wget \
    git 

# install Go 1.8
ENV GOPATH /root/go
ENV GO_VERSION 1.8
ENV GO_ARCH amd64
RUN wget https://storage.googleapis.com/golang/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz; \
   tar -C /usr/local/ -xf /go${GO_VERSION}.linux-${GO_ARCH}.tar.gz ; \
   rm /go${GO_VERSION}.linux-${GO_ARCH}.tar.gz

ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH

WORKDIR /
