PACKAGES=$(shell find . -name '*.go' -print0 | xargs -0 -n1 dirname | sort --unique)
GOFILES= vflow.go ipfix.go sflow.go options.go stats.go 
LDFLAGS= -ldflags "-X main.version=0.3.0"

default: test

test:
	go test -v ./... -timeout 1m

bench:
	go test -v ./... -bench=. -timeout 2m

run:
	cd vflow; go run $(GOFILES) -sflow-workers 100 -ipfix-workers 100

debug:
	cd vflow; go run $(GOFILES) -sflow-workers 100 -ipfix-workers 100 -verbose=true

gctrace:
	cd vflow; env GODEBUG=gctrace=1 go run $(GOFILES) -sflow-workers 100 -ipfix-workers 100

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
	cd vflow; go build $(LDFLAGS) -o vflow $(GOFILES)
