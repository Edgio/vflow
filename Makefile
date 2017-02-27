PACKAGES=$(shell find . -name '*.go' -print0 | xargs -0 -n1 dirname | sort --unique)
default:
	go test -v ./...
test:
	go test -v ./...

bench:
	go test -v ./... -bench=.

run:
	cd vflow; go run *.go -sflow-workers 100 -ipfix-workers 100 -verbose=false

debug:
	cd vflow; go run *.go -sflow-workers 100 -ipfix-workers 100 -verbose=true

gctrace:
	cd vflow; env GODEBUG=gctrace=1 go run *.go -sflow-workers 100 -ipfix-workers 100 -verbose=false

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
build:
	cd vflow; go build -o vflow *.go 
