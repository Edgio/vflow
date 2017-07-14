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
	sed -i 's/%VERSION%/${VERSION}/' ${DEBPATH}/DEBIAN/control
	mkdir -p ${DEBPATH}/usr/local/bin
	cp vflow/vflow ${DEBPATH}/usr/local/bin
	dpkg-deb -b ${DEBPATH}
	mv ${DEBPATH}.deb scripts/vflow${VERSION}.deb
	sed -i 's/${VERSION}/%VERSION%/' ${DEBPATH}/DEBIAN/control
