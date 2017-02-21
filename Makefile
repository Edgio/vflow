PACKAGES=$(shell find . -name '*.go' -print0 | xargs -0 -n1 dirname | sort --unique)
default:

test:
	go test ./...

bench:
	go test ./... -bench=.

run:
	cd vflow; go run options.go sflow.go ipfix.go stats.go vflow.go -sflow-workers 100 -ipfix-workers 100 -verbose=false

debug:
	cd vflow; go run options.go sflow.go ipfix.go stats.go vflow.go -sflow-workers 100 -ipfix-workers 100 -verbose=true

gctrace:
	cd vflow; env GODEBUG=gctrace=1 go run options.go sflow.go ipfix.go stats.go vflow.go -sflow-workers 100 -ipfix-workers 100 -verbose=false

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
