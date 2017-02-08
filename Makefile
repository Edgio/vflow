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

