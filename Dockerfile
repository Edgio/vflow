FROM ubuntu:trusty

ENV GOPATH /root/go
ENV GO_VERSION 1.8
ENV GO_ARCH amd64
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    wget \
    git \
    supervisor

# install Go 1.8
RUN wget -q https://storage.googleapis.com/golang/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz; \
   tar -C /usr/local/ -xf /go${GO_VERSION}.linux-${GO_ARCH}.tar.gz ; \
   rm /go${GO_VERSION}.linux-${GO_ARCH}.tar.gz

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin && chmod -R 777 ${GOPATH}

# build vFlow
RUN mkdir -p ${GOPATH}/src/github.com/VerizonDigital/vflow
ADD . ${GOPATH}/src/github.com/VerizonDigital/vflow
WORKDIR ${GOPATH}/src/github.com/VerizonDigital/vflow/vflow
RUN ["go", "get", "-d", "./..."]
RUN ["go", "build", "-o", "/usr/local/bin/vflow"]

ADD scripts/vflow.supervisor /etc/supervisor/conf.d/vflow.conf

EXPOSE 4739 6343 4729 8081

CMD ["supervisord", "-n"]
