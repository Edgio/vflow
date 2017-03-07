PACKAGES=$(shell find . -name '*.go' -print0 | xargs -0 -n1 dirname | sort --unique)
GOFILES= vflow.go ipfix.go sflow.go options.go stats.go 

default:
	go test -v ./... -timeout 1m
test:
	go test -v ./... -timeout 1m

bench:
	go test -v ./... -bench=. -timeout 2m

run:
	cd vflow; go run $(GOFILES) -sflow-workers 100 -ipfix-workers 100 -verbose=false

debug:
	cd vflow; go run $(GOFILES) -sflow-workers 100 -ipfix-workers 100 -verbose=true

gctrace:
	cd vflow; env GODEBUG=gctrace=1 go run $(GOFILES) -sflow-workers 100 -ipfix-workers 100 -verbose=false

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

build:
	cd vflow; go build -o vflow $(GOFILES)
