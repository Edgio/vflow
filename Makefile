PACKAGES=$(shell find . -name '*.go' -print0 | xargs -0 -n1 dirname | sort --unique)
GOFILES= vflow.go ipfix.go sflow.go netflow_v9.go options.go stats.go 
LDFLAGS= -ldflags "-X main.version=0.3.1"

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

build-windows: depends
	cd vflow; gox $(LDFLAGS) 


vflow/vflow_windows_386.exe: depends
	cd vflow; gox $(LDFLAGS) -os="windows" -arch="386"

vflow/vflow_windows_amd64.exe: depends
	cd vflow; gox $(LDFLAGS) -os="windows" -arch="amd64"

vflow/vflow_darwin_386: depends
	cd vflow; gox $(LDFLAGS) -os="darwin" -arch="386"

vflow/vflow_darwin_amd64: depends
	cd vflow; gox $(LDFLAGS) -os="darwin" -arch="amd64"

vflow/vflow_freebsd_386: depends
	cd vflow; gox $(LDFLAGS) -os="freebsd" -arch="386"

vflow/vflow_freebsd_amd64: depends
	cd vflow; gox $(LDFLAGS) -os="freebsd" -arch="amd64"

vflow/vflow_linux_386: depends
	cd vflow; gox $(LDFLAGS) -os="linux" -arch="386"

vflow/vflow_linux_amd64: depends
	cd vflow; gox $(LDFLAGS) -os="linux" -arch="amd64"

cross-compile: vflow/vflow_windows_386.exe vflow/vflow_windows_amd64.exe vflow/vflow_darwin_386 vflow/vflow_darwin_amd64 vflow/vflow_freebsd_386 vflow/vflow_freebsd_amd64 vflow/vflow_linux_386 vflow/vflow_linux_amd64
