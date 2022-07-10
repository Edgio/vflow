## Command line Kafka consumer without any back-end.
This is an example of the vFlow.IPFIX CLI consumer. it supports only destination IP address filtring by default but you can change the element ID number through CLI based on the IANA IPFIX element ID.

### Build

```
go get -d ./...
go build main.go
```
