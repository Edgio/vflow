VERSION= 0.3.2
PACKAGES= $(shell find . -name '*.go' -print0 | xargs -0 -n1 dirname | sort --unique)
LDFLAGS= -ldflags "-X main.version=${VERSION}"
DEBPATH= scripts/dpkg

default: test

test:
	go test -v ./... -timeout 1m

bench:
	go test -v ./... -bench=. -timeout 2m

run: build
	cd vflow; ./vflow -sflow-workers 100 -ipfix-workers 100

debug: build
	cd vflow; ./vflow -sflow-workers 100 -ipfix-workers 100 -verbose=true

gctrace: build
	cd vflow; env GODEBUG=gctrace=1 ./vflow -sflow-workers 100 -ipfix-workers 100

lint:
	golint ./...

cyclo:
	gocyclo -over 15 $(PACKAGES)

errcheck:
	errcheck ./...

tools:
	go get github.com/golang/lint/golint
	go get github.com/kisielk/errcheck
	go get github.com/alecthomas/gocyclo

depends:
	go get -d ./...

build: depends
	cd vflow; go build $(LDFLAGS)

dpkg: build
	mkdir -p ${DEBPATH}/etc/init.d ${DEBPATH}/etc/logrotate.d
	mkdir -p ${DEBPATH}/etc/vflow ${DEBPATH}/usr/share/doc/vflow
	mkdir -p ${DEBPATH}/usr/bin ${DEBPATH}/usr/local/vflow
	sed -i 's/%VERSION%/${VERSION}/' ${DEBPATH}/DEBIAN/control
	cp vflow/vflow ${DEBPATH}/usr/bin/
	cp scripts/vflow.service ${DEBPATH}/etc/init.d/vflow
	cp scripts/vflow.logrotate ${DEBPATH}/etc/logrotate.d/vflow
	cp scripts/vflow.conf ${DEBPATH}/etc/vflow/vflow.conf
	cp scripts/kafka.conf ${DEBPATH}/etc/vflow/mq.conf
	cp ${DEBPATH}/DEBIAN/copyright ${DEBPATH}/usr/share/doc/vflow/
	dpkg-deb -b ${DEBPATH}
	mv ${DEBPATH}.deb scripts/vflow${VERSION}.deb
	sed -i 's/${VERSION}/%VERSION%/' ${DEBPATH}/DEBIAN/control
